package service

import (
	"context"
	"encoding/json"
	errs "errors"
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
	mq amqp.AmqpClient
	rs res.Store
	cs chat.Store
}

func NewChatService(us auth.Store, gs group.Store, rs res.Store, mq amqp.AmqpClient, cs chat.Store) *ChatService {
	return &ChatService{
		us: us,
		gs: gs,
		mq: mq,
		rs: rs,
		cs: cs,
	}
}

var _ chat.Service = &ChatService{}

// CreateUserExchange will create the AMQP exchange to receive user messages
// This exchange will be bound to queues representing different Websocket clients
// for the same user (if a user is using multiple devices to connect, he
// will get Websocket notifications on all devices)
func (c ChatService) CreateUserExchange(ctx context.Context, userKey model.UserKey) (string, error) {

	ctx, l := GetCtx(ctx, "ChatService", "CreateUserExchange")

	amqpChannel, err := c.mq.GetChannel()
	if err != nil {
		l.Error("could not get amqp channel", zap.Error(err))
		return "", err
	}
	defer amqpChannel.Close()

	exchangeName := c.GetUserExchangeName(ctx, userKey)

	err = amqpChannel.ExchangeDeclare(ctx, exchangeName, "fanout", true, false, false, false, nil)
	if err != nil {
		l.Error("could not declare user exchange", zap.Error(err))
		return "", err
	}

	return exchangeName, nil
}

// GetUserExchangeName will return the name of the exchange for a user key
func (c ChatService) GetUserExchangeName(ctx context.Context, userKey model.UserKey) string {
	return "users." + userKey.String()
}

// CreateChannel Will create a channel by key
func (c ChatService) CreateChannel(ctx context.Context, channelKey model.ChannelKey, channelType chat.ChannelType) (*chat.Channel, error) {

	ctx, l := GetCtx(ctx, "ChatService", "CreateChannel")

	channel := &chat.Channel{
		ID:    channelKey.ID,
		Title: "",
		Type:  channelType,
	}

	err := c.cs.CreateChannel(ctx, channel)
	if err != nil {
		l.Error("could not create channel", zap.Error(err))
		return nil, err
	}

	channel, err = c.cs.GetChannel(ctx, channelKey)
	if err != nil {
		l.Error("could not get channel", zap.Error(err))
		return nil, err
	}

	return channel, nil
}

// SubscribeToChannel will subscribe a user to a given channel
func (c ChatService) SubscribeToChannel(ctx context.Context, channelSubscriptionKey model.ChannelSubscriptionKey, name string) (*chat.ChannelSubscription, error) {

	ctx, l := GetCtx(ctx, "ChatService", "SubscribeToChannel")
	l = l.With(zap.Object("channel_subscription", channelSubscriptionKey))

	l.Debug("subscribing to channel")

	channelSubscription, err := c.cs.CreateSubscription(ctx, channelSubscriptionKey, name)
	if err != nil {
		l.Error("could not create channel subscription", zap.Error(err))
		return nil, err
	}

	amqpChannel, err := c.mq.GetChannel()
	if err != nil {
		l.Error("could not get amqp channel", zap.Error(err))
		return nil, err
	}
	defer amqpChannel.Close()

	userExchangeName, err := c.CreateUserExchange(ctx, channelSubscriptionKey.UserKey)
	if err != nil {
		return nil, err
	}

	headers := c.getChannelBindingHeaders(channelSubscriptionKey)
	err = amqpChannel.ExchangeBind(ctx, userExchangeName, "", amqp.WebsocketMessagesExchange, false, headers)
	if err != nil {
		l.Error("could not bind user exchange", zap.Error(err))
		return nil, err
	}

	return channelSubscription, nil

}

func (c ChatService) UnsubscribeFromChannel(ctx context.Context, channelSubscriptionKey model.ChannelSubscriptionKey) error {

	ctx, l := GetCtx(ctx, "ChatService", "UnsubscribeFromChannel")
	l = l.With(zap.Object("channel_subscription", channelSubscriptionKey))

	l.Debug("unsubscribing from channel")

	err := c.cs.DeleteSubscription(ctx, channelSubscriptionKey)
	if err != nil {
		l.Error("could not delete channel subscription", zap.Error(err))
		return err
	}

	amqpChannel, err := c.mq.GetChannel()
	if err != nil {
		l.Error("could not get amqp channel", zap.Error(err))
		return err
	}
	defer amqpChannel.Close()

	userExchangeName := channelSubscriptionKey.UserKey.GetExchangeName()
	headers := c.getChannelBindingHeaders(channelSubscriptionKey)
	err = amqpChannel.ExchangeUnbind(ctx, userExchangeName, "", amqp.WebsocketMessagesExchange, false, headers)
	if err != nil {
		l.Error("could not bind user exchange", zap.Error(err))
		return err
	}

	return nil

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
						c.getUserLink(loggedInUser),
						c.getResourceLink(resource),
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

// getUserLink Gets the markdown representing the link to a user profile
func (c ChatService) getUserLink(loggedInUser model.UserReference) string {
	userKey := loggedInUser.GetUserKey()
	userLink := fmt.Sprintf("[/users/%s](%s)", userKey.String(), loggedInUser.GetUsername())
	return userLink
}

// getResourceLink Gets the markdown representing the link to a resource
func (c ChatService) getResourceLink(resource *res.Resource) string {
	userId := resource.CreatedBy
	needOrOffer := "offers"
	if resource.Type == res.ResourceRequest {
		needOrOffer = "needs"
	}
	resourceId := resource.ID.String()
	return fmt.Sprintf("[/users/%s/%s/%s](%s)", userId, needOrOffer, resourceId, resource.Summary)
}

func (c ChatService) getUserSessionAndKey(ctx context.Context) (*auth.UserSession, model.UserKey, error) {
	userSession, err := auth.GetUserSession(ctx)
	if err != nil {
		return nil, model.UserKey{}, err
	}
	return userSession, userSession.GetUserKey(), nil
}

// getOrCreateConversationChannel Will retrieve or create the conversation channel for a given set of users
// If the channel is not already created, it will automatically be created and the users will be subscribed to it.
// It also sets up RabbitMQ routing so that messages to this channel will find the user's exchange.
func (c ChatService) getOrCreateConversationChannel(ctx context.Context, userKeys *model.UserKeys) (*chat.GetOrCreateConversationChannelResponse, error) {

	ctx, l := GetCtx(ctx, "ChatService", "getOrCreateConversationChannel")


	// Retrieve the channel key for that conversation
	channelKey, err := c.getConversationChannelKey(ctx, userKeys)
	if err != nil {
		l.Error("could not get conversation channel key", zap.Error(err))
		return nil, err
	}


	channel, err := c.cs.GetChannel(nil, channelKey)
	if err != nil && !errs.Is(err, chat.ErrChannelNotFound) {
		l.Error("could not get channel", zap.Error(err))
		return nil, err
	}

	// If the channel doesn't exist
	if errs.Is(err, chat.ErrChannelNotFound) {

		l.Debug("channel doesn't exist. Creating channel")

		var title = ""

		newChannel := &chat.Channel{
			ID:    channelKey.String(),
			Title: title,
			Type:  chat.ConversationChannel,
		}

		err := c.cs.CreateChannel(ctx, newChannel)
		if err != nil {
			l.Error("could not create channel", zap.Error(err))
			return nil, err
		}

		amqpChannelKey, err := c.getConversationChannelKey(ctx, userKeys)
		if err != nil {
			l.Error("could not get conversation channel key")
			return nil, err
		}

		_, err = c.createSubscriptionsAndMqBindingsForUserConversation(ctx, userKeys)
		if err != nil {
			l.Error("could not create subscriptions and mq bindings", zap.Error(err))
			return nil, err
		}

		channel, err = c.cs.GetChannel(ctx, amqpChannelKey)
		if err != nil {
			l.Error("could not get channel", zap.Error(err))
			return nil, err
		}

	}

	return &chat.GetOrCreateConversationChannelResponse{
		Channel: channel,
	}, nil

}

func (c ChatService) SendChannelMessage(ctx context.Context, channelKey model.ChannelKey, message string) (*chat.Message, error) {

	ctx, l := GetCtx(ctx, "ChatService", "SendConversationMessage")

	userSession, err := auth.GetUserSession(ctx)
	if err != nil {
		l.Error("cannot get user session", zap.Error(err))
		return nil, err
	}

	_, err = c.cs.GetChannel(ctx, channelKey)
	if err != nil {
		l.Error("cannot get channel", zap.Error(err))
		return nil, err
	}

	authUserKey := userSession.GetUserKey()
	sentMessage, err := c.cs.SaveMessage(ctx, chat.NewSaveMessageRequest(
		channelKey,
		authUserKey,
		userSession.GetUsername(),
		message,
		[]chat.Block{},
		[]chat.Attachment{},
		nil,
	))


	amqpChannel, err := c.mq.GetChannel()
	if err != nil {
		l.Error("could not get amqp channel", zap.Error(err))
		return nil, err
	}
	defer amqpChannel.Close()


	evt := amqp.Event{
		Type:      "message",
		SubType:   "user",
		Channel:   channelKey.String(),
		User:      authUserKey.String(),
		ID:        sentMessage.Message.Key.String(),
		Timestamp: sentMessage.Message.SentAt.String(),
		Text:      sentMessage.Message.Text,
	}

	js, err := json.Marshal(evt)
	if err != nil {
		l.Error("could not marshal message", zap.Error(err))
		return nil, err
	}

	err = amqpChannel.Publish(ctx, amqp.MessagesExchange, "", false, false, amqp.AmqpPublishing{
		Headers: map[string]interface{}{
			"channel_id": channelKey.String(),
			"event_type": "chat.message",
		},
		ContentType: "application/json",
		Body:        js,
	})

	if err != nil {
		l.Error("failed to publish message", zap.Error(err))
		return nil, err
	}

	return nil, nil
}

func (c ChatService) SendConversationMessage(ctx context.Context, request *chat.SendConversationMessage) (*chat.SendConversationMessageResponse, error) {

	ctx, l := GetCtx(ctx, "ChatService", "SendConversationMessage")

	createdChannel, err := c.getOrCreateConversationChannel(ctx, request.ToUserKeys)
	if err != nil {
		l.Error("could not get or create conversation channel", zap.Error(err))
		return nil, err
	}

	channelKey := createdChannel.Channel.GetKey()
	sendMessageRequest := chat.NewSaveMessageRequest(
		channelKey,
		request.FromUserKey,
		request.FromUserName,
		request.Text,
		request.Blocks,
		request.Attachments,
		request.OnlyVisibleToUserKey,
	)
	sentMessage, err := c.cs.SaveMessage(ctx, sendMessageRequest)

	mjs, _ := json.Marshal(sendMessageRequest)
	fmt.Println(string(mjs))


	if err != nil {
		l.Error("could not save message", zap.Error(err))
		return nil, err
	}

	amqpChannel, err := c.mq.GetChannel()
	if err != nil {
		l.Error("could not get amqp channel", zap.Error(err))
		return nil, err
	}
	defer amqpChannel.Close()

	evt := amqp.Event{
		Type:      "message",
		SubType:   "user",
		Channel:   channelKey.String(),
		User:      request.FromUserKey.String(),
		ID:        sentMessage.Message.Key.String(),
		Timestamp: sentMessage.Message.SentAt.String(),
		Text:      sentMessage.Message.Text,
	}

	js, err := json.Marshal(evt)
	if err != nil {
		l.Error("could not marshal message", zap.Error(err))
		return nil, err
	}

	err = amqpChannel.Publish(ctx, amqp.MessagesExchange, "", false, false, amqp.AmqpPublishing{
		Headers: map[string]interface{}{
			"channel_id": channelKey.String(),
			"event_type": "chat.message",
		},
		ContentType: "application/json",
		Body:        js,
	})

	if err != nil {
		l.Error("failed to publish message", zap.Error(err))
		return nil, err
	}

	return &chat.SendConversationMessageResponse{
		Message: sentMessage.Message,
	}, nil
}

func (c ChatService) SendGroupMessage(ctx context.Context, request *chat.SendGroupMessage) (*chat.SendGroupMessageResponse, error) {

	ctx, l := GetCtx(ctx, "ChatService", "SendConversationMessage")

	channelKey := request.GroupKey.GetChannelKey()
	sentMessage, err := c.cs.SaveMessage(ctx, chat.NewSaveMessageRequest(
		channelKey,
		request.FromUserKey,
		request.FromUserName,
		request.Text,
		request.Blocks,
		request.Attachments,
		request.OnlyVisibleToUserKey,
	))

	if err != nil {
		l.Error("could not save message", zap.Error(err))
		return nil, err
	}

	l.Debug("getting amqp channel")

	amqpChannel, err := c.mq.GetChannel()
	if err != nil {
		l.Error("could not get amqp channel", zap.Error(err))
		return nil, err
	}
	defer amqpChannel.Close()

	l.Debug("sending message to RabbitMQ")

	evt := amqp.Event{
		Type:      "message",
		SubType:   "user",
		Channel:   channelKey.String(),
		User:      request.FromUserKey.String(),
		ID:        sentMessage.Message.Key.String(),
		Timestamp: sentMessage.Message.SentAt.String(),
		Text:      sentMessage.Message.Text,
	}

	js, err := json.Marshal(evt)
	if err != nil {
		l.Error("could not marshal message", zap.Error(err))
		return nil, err
	}

	err = amqpChannel.Publish(ctx, amqp.MessagesExchange, "", false, false, amqp.AmqpPublishing{
		Headers: map[string]interface{}{
			"channel_id": channelKey.String(),
			"event_type": "chat.message",
		},
		ContentType: "application/json",
		Body:        js,
	})

	if err != nil {
		l.Error("failed to publish message", zap.Error(err))
		return nil, err
	}

	// todo
	return &chat.SendGroupMessageResponse{}, nil

}

// createSubscriptionsAndMqBindingsForUserConversation subscribes multiple users to a conversation channel
// It also generates the RabbitMQ bindings so that messages to the conversation
// are also properly routed to the user's exchange.
func (c ChatService) createSubscriptionsAndMqBindingsForUserConversation(ctx context.Context, userKeys *model.UserKeys) ([]chat.ChannelSubscription, error) {

	ctx, l := GetCtx(ctx, "ChatService", "createSubscriptionsAndMqBindingsForUserConversation")

	channelKey, err := c.getConversationChannelKey(ctx, userKeys)
	if err != nil {
		l.Error("could not get conversation channel key")
		return nil, err
	}

	users, err := c.us.GetByKeys(ctx, userKeys.Items)
	if err != nil {
		l.Error("could not get users by keys")
		return nil, err
	}

	var subscriptions []chat.ChannelSubscription
	for _, user := range users.Items {
		subscription, err := c.createSubscriptionAndMqBindingForUserConversation(ctx, user, users, channelKey)
		if err != nil {
			l.Error("could not create subscription and mq binding for user conversation", zap.Error(err))
			return nil, err
		}
		subscriptions = append(subscriptions, *subscription)
	}
	return subscriptions, nil
}

// createSubscriptionAndMqBindingForUserConversation subscribes the user to a conversation channel.
// It also generates the RabbitMQ bindings so that messages to the conversation
// are also properly routed to the user's exchange.
func (c ChatService) createSubscriptionAndMqBindingForUserConversation(
	ctx context.Context,
	user auth.User,
	conversationUsers auth.Users,
	channelKey model.ChannelKey,
) (*chat.ChannelSubscription, error) {

	ctx, l := GetCtx(ctx, "ChatService", "createSubscriptionAndMqBindingForUserConversation")

	// The name of the conversation
	conversationName := c.getConversationNameForUser(ctx, conversationUsers, user)

	channelSubscriptionKey := model.NewChannelSubscriptionKey(channelKey, user.GetUserKey())
	subscription, err := c.cs.CreateSubscription(ctx, channelSubscriptionKey, conversationName)
	if err != nil {
		l.Error("could not create subscription")
		return nil, err
	}

	amqpChannel, err := c.mq.GetChannel()
	if err != nil {
		l.Error("could not get amqp channel", zap.Error(err))
		return nil, err
	}
	defer amqpChannel.Close()

	headers := c.getChannelBindingHeaders(subscription.GetKey())
	userExchangeName, err := c.CreateUserExchange(ctx, user.GetUserKey())
	err = amqpChannel.ExchangeBind(ctx, userExchangeName, "", amqp.WebsocketMessagesExchange, false, headers)
	if err != nil {
		l.Error("could not register user channel binding", zap.Error(err))
		return nil, err
	}

	return subscription, nil
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
	sort.Slice(userList, func(i, j int) bool {
		return userList[i].Username > userList[j].Username
	})

	// Get the usernames, but omit the user that will see this name
	var userNames []string
	for _, otherUser := range userList {
		if u.GetUserKey() == otherUser.GetUserKey() {
			continue
		}
		userNames = append(userNames, otherUser.Username)
	}

	// Join the name with a space
	conversationName := strings.Join(userNames, " ")

	return conversationName
}

// getConversationChannelKey Returns the id of a conversation between users
// Only a single conversation can exist between a group of people.
// There can only be one conversation with Joe, Dana and Mark.
// So the ID of the conversation is composed of the
// sorted IDs of its participants.
func (c ChatService) getConversationChannelKey(ctx context.Context, participants *model.UserKeys) (model.ChannelKey, error) {

	ctx, l := GetCtx(ctx, "ChatService", "getConversationChannelKey")

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
