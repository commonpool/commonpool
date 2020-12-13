package store

import (
	"context"
	"encoding/json"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/store"
	"gorm.io/gorm"
)

type ChatStore struct {
	db *gorm.DB
	as auth.Store
	mq amqp.Client
}

var _ chat.Store = &ChatStore{}

func NewChatStore(db *gorm.DB, as auth.Store, amqpClient amqp.Client) *ChatStore {
	return &ChatStore{
		as: as,
		db: db,
		mq: amqpClient,
	}
}

//
// func (cs *ChatStore) GetOrCreateConversationChannel(request *chat.GetOrCreateConversationChannel) (*chat.GetOrCreateConversationChannelResponse, error) {
//
// 	var channelKey *model.ConversationKey
// 	var err error
//
// 	if request.Type == chat.ConversationChannel {
// 		channelKey, err = cs.getConversationId(request.ParticipantList)
// 	} else if request.Type == chat.GroupChannel {
// 		groupKey := request.Group.GetKey()
// 		channelKey, err = cs.getGroupChannelKey(&groupKey)
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	channel, err := cs.getChannel(channelKey)
//
// 	if errors.Is(err, gorm.ErrRecordNotFound) {
//
// 		var title = ""
// 		if request.Type == chat.GroupChannel {
// 			title = request.Group.Name
// 		}
//
// 		channel = &chat.Channel{
// 			ID:    channelKey.String(),
// 			Title: title,
// 			Type:  chat.ConversationChannel,
// 		}
//
// 		err := cs.db.
// 			Model(chat.Channel{}).
// 			Create(channel).
// 			Error
//
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		if request.Type == chat.ConversationChannel {
// 			err = cs.createConversationSubscriptions(request, channelKey)
// 			if err != nil {
// 				return nil, err
// 			}
// 		} else if request.Type == chat.GroupChannel {
// 			err := cs.createGroupSubscriptions(request, channelKey)
// 			if err != nil {
// 				return nil, err
// 			}
// 		}
//
// 		channel, err = cs.getChannel(channelKey)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 	}
//
// 	return &chat.GetOrCreateConversationChannelResponse{
// 		Channel:         channel,
// 		ParticipantList: request.ParticipantList,
// 	}, nil
//
// }

// func (cs *ChatStore) createGroupSubscriptions(request *chat.GetOrCreateConversationChannel, key *model.ConversationKey) error {
// 	panic("implement me")
// }
//
// func (cs *ChatStore) createConversationSubscriptions(request *chat.GetOrCreateConversationChannel, channelKey *model.ConversationKey) error {
// 	users, err := cs.as.GetByKeys(nil, request.ParticipantList.Items)
// 	if err != nil {
// 		return err
// 	}
// 	userList := users.Items
// 	sort.Slice(userList, func(i, j int) bool {
// 		return users.Items[i].Username > users.Items[j].Username
// 	})
//
// 	var subscriptions []chat.ChannelSubscription
// 	for _, key := range request.ParticipantList.Items {
//
// 		var userNames []string
// 		for _, user := range userList {
// 			if user.GetUserKey() == key {
// 				continue
// 			}
// 			userNames = append(userNames, user.Username)
// 		}
//
// 		subscription := chat.ChannelSubscription{
// 			Name:          strings.Join(userNames, " "),
// 			UserID:        key.String(),
// 			ChannelID:     channelKey.String(),
// 			LastMessageAt: time.Unix(0, 0),
// 			LastTimeRead:  time.Unix(0, 0),
// 		}
// 		subscriptions = append(subscriptions, subscription)
//
// 		err := cs.mq.BindUserExchangeToChannel(nil, channelKey.String(), key.String())
// 		if err != nil {
// 			return err
// 		}
// 	}
//
// 	err = cs.db.Model(chat.ChannelSubscription{}).Create(subscriptions).Error
// 	if err != nil {
// 		for _, userKey := range request.ParticipantList.Items {
// 			_ = cs.mq.UnbindUserExchangeFromChannel(channelKey.String(), userKey.String())
// 		}
// 		return err
// 	}
//
// 	return nil
// }

func mapMessage(ctx context.Context, message *store.Message) (*chat.Message, error) {

	var blocks []chat.Block
	var attachments []chat.Attachment
	if err := json.Unmarshal([]byte(message.Blocks), &blocks); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(message.Attachments), &attachments); err != nil {
		return nil, err
	}

	var visibleToUser *model.UserKey
	if message.VisibleToUser != nil {
		visibleToUserKey := model.NewUserKey(*message.VisibleToUser)
		visibleToUser = &visibleToUserKey
	}

	returnMessage := &chat.Message{
		Key:            model.NewMessageKey(message.ID),
		ChannelKey:     model.NewConversationKey(message.ChannelID),
		MessageType:    message.MessageType,
		MessageSubType: message.MessageSubType,
		SentBy: chat.MessageSender{
			Type:     chat.UserMessageSender,
			UserKey:  model.NewUserKey(message.SentById),
			USername: message.SentByUsername,
		},
		SentAt:        message.SentAt,
		Text:          message.Text,
		Blocks:        blocks,
		Attachments:   attachments,
		VisibleToUser: visibleToUser,
	}
	return returnMessage, nil
}
