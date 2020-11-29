package service

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	res "github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/utils"
	"go.uber.org/zap"
	"sort"
	"strings"
)

type ChatService struct {
	us auth.Store
	gs group.Store
	mq amqp.Client
	rs res.Store
	cs chat.Store
}

func NewChatService(us auth.Store, gs group.Store, rs res.Store, mq amqp.Client, cs chat.Store) *ChatService {
	return &ChatService{
		us: us,
		gs: gs,
		mq: mq,
		rs: rs,
		cs: cs,
	}
}

var _ chat.Service = &ChatService{}

// GetUserExchangeName will return the name of the exchange for a user key
func (c ChatService) GetUserExchangeName(ctx context.Context, userKey model.UserKey) string {
	return "users." + userKey.String()
}

// getChannelBindingHeaders returns the binding headers to link the websocket messages exchange and a given user exchange
// The user will receive messages on his exchange if the message has a "channel_id" = "subscribed channel id" header
func (c ChatService) getChannelBindingHeaders(channelSubscriptionKey model.ChannelSubscriptionKey) map[string]interface{} {
	return map[string]interface{}{
		"event_type": "chat.message",
		"channel_id": channelSubscriptionKey.ChannelKey.String(),
		"x-match":    "all",
	}
}

func (c ChatService) DeleteGroupChannel(ctx context.Context, request *chat.DeleteGroupChannel) (*chat.DeleteGroupChannelResponse, error) {
	panic("implement me")
}

// NotifyUserInterestedAboutResource will create a channel between two users if it doesn't exist,
// and will send a message to the owner of the resource notifying them that someone is interested.
func (c ChatService) NotifyUserInterestedAboutResource(ctx context.Context, request *chat.NotifyUserInterestedAboutResource) (*chat.NotifyUserInterestedAboutResourceResponse, error) {

	ctx, l := GetCtx(ctx, "ChatService", "NotifyUserInterestedAboutResource")

	l.Debug("getting logged in user")

	loggedInUser, loggedInUserKey, err := c.getUserSessionAndKey(ctx)
	if err != nil {
		return nil, err
	}

	l.Debug("retrieving resource")

	getResource := c.rs.GetByKey(ctx, res.NewGetResourceByKeyQuery(request.ResourceKey))
	if getResource.Error != nil {
		l.Error("could not get resource", zap.Error(err))
		return nil, getResource.Error
	}
	resource := getResource.Resource
	resourceOwnerKey := resource.GetOwnerKey()

	l.Debug("make sure resource owner is not inquiring about his own resource")

	// make sure auth user is not resource owner
	// doesn't make sense for one to inquire about his own stuff
	if resourceOwnerKey == loggedInUserKey {
		err := errors.ErrCannotInquireAboutOwnResource()
		l.Error("user cannot inquire about his own resource", zap.Error(err))
		return nil, err
	}

	userKeys := model.NewUserKeys([]model.UserKey{loggedInUserKey, resourceOwnerKey})

	l.Debug("sending private message to resource owner")

	_, err = c.SendConversationMessage(ctx, chat.NewSendConversationMessage(
		loggedInUserKey,
		loggedInUser.Username,
		userKeys,
		request.Message,
		[]chat.Block{
			*chat.NewHeaderBlock(chat.NewMarkdownObject("Someone is interested in your stuff!"), nil),
			*chat.NewContextBlock([]chat.BlockElement{
				chat.NewMarkdownObject(
					fmt.Sprintf("%s is interested by your post %s.",
						c.GetUserLink(loggedInUserKey),
						c.GetResourceLink(request.ResourceKey),
					),
				),
			}, nil),
		},
		[]chat.Attachment{},
		&resourceOwnerKey,
	))
	if err != nil {
		l.Error("could not send message to resource owner", zap.Error(err))
		return nil, err
	}

	l.Debug("sending message to conversation")

	sentPublicMessage, err := c.SendConversationMessage(ctx, chat.NewSendConversationMessage(
		loggedInUserKey,
		loggedInUser.Username,
		userKeys,
		request.Message,
		[]chat.Block{},
		[]chat.Attachment{},
		nil,
	))

	if err != nil {
		l.Error("could not send message to conversation", zap.Error(err))
		return nil, err
	}

	return &chat.NotifyUserInterestedAboutResourceResponse{
		ChannelKey: sentPublicMessage.Message.ChannelKey,
	}, nil

}

// GetUserLink Gets the markdown representing the link to a user profile
func (c ChatService) GetUserLink(userKey model.UserKey) string {
	return fmt.Sprintf("<commonpool-user id='%s'></commonpool-user>", userKey.String())
}

// GetResourceLink Gets the markdown representing the link to a resource
func (c ChatService) GetResourceLink(resource model.ResourceKey) string {
	return fmt.Sprintf("<commonpool-resource id='%s'><commonpool-resource>", resource.String())
}

func (c ChatService) getUserSessionAndKey(ctx context.Context) (*auth.UserSession, model.UserKey, error) {
	userSession, err := auth.GetUserSession(ctx)
	if err != nil {
		return nil, model.UserKey{}, err
	}
	return userSession, userSession.GetUserKey(), nil
}

// getConversationNameForUser Gets the name of the conversation for a specific user.
// A user will see the name of a conversation as equal to the names of the other participants
// but not his own name in it.
//
// Joe  would see "Dana, Mark"
// Dana would see "Joe, Mark"
// Mark would see "Dana, Joe"
func (c ChatService) getConversationNameForUser(
	ctx context.Context,
	us auth.Users,
	u auth.User,
) string {

	// First, sort the user names
	userList := us.Items

	copied := make([]auth.User, len(userList))
	copy(copied, userList)
	sort.Slice(copied, func(i, j int) bool {
		return copied[i].Username > copied[j].Username
	})

	// Get the usernames, but omit the user that will see this name
	var userNames []string
	for _, otherUser := range copied {
		if u.GetUserKey() == otherUser.GetUserKey() {
			continue
		}
		userNames = append(userNames, otherUser.Username)
	}

	// Join the name with a space
	conversationName := strings.Join(userNames, " ")

	return conversationName
}

// GetConversationChannelKey Returns the id of a conversation between users
// Only a single conversation can exist between a group of people.
// There can only be one conversation with Joe, Dana and Mark.
// So the ID of the conversation is composed of the
// sorted IDs of its participants.
func (c ChatService) GetConversationChannelKey(ctx context.Context, participants *model.UserKeys) (model.ChannelKey, error) {

	ctx, l := GetCtx(ctx, "ChatService", "GetConversationChannelKey")

	if participants == nil || len(participants.Items) == 0 {
		err := fmt.Errorf("cannot get conversation channel for 0 participants")
		l.Error(err.Error())
		return model.ChannelKey{}, err
	}

	var shortUids []string
	for _, participant := range participants.Items {
		sid, err := utils.ShortUuidFromStr(participant.String())
		if err != nil {
			return model.ChannelKey{}, err
		}
		shortUids = append(shortUids, sid)
	}
	sort.Strings(shortUids)
	channelId := strings.Join(shortUids, "-")
	channelKey := model.NewConversationKey(channelId)

	return channelKey, nil
}
