package store

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/logging"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/utils"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sort"
	"strings"
	"time"
)

type ChatStore struct {
	db *gorm.DB
	as auth.Store
	mq amqp.Client
}

func (cs *ChatStore) DeleteSubscription(ctx context.Context, key model.ChannelSubscriptionKey) error {

	ctx, l := GetCtx(ctx, "ChatStore", "CreateSubscription")
	l = l.With(zap.Object("channel_subscription", key))

	err := cs.db.Delete(chat.ChannelSubscription{},
		"user_id = ? and channel_id = ?",
		key.UserKey.String(),
		key.ChannelKey.String()).
		Error

	if err != nil {
		l.Error("could not delete channel subscription", zap.Error(err))
		return err
	}

	return nil
}

func (cs *ChatStore) CreateChannel(ctx context.Context, channel *chat.Channel) error {
	return cs.db.WithContext(ctx).Create(channel).Error
}

var _ chat.Store = &ChatStore{}

func NewChatStore(db *gorm.DB, as auth.Store, amqpClient amqp.Client) *ChatStore {
	return &ChatStore{
		as: as,
		db: db,
		mq: amqpClient,
	}
}

func (cs *ChatStore) CreateSubscription(ctx context.Context, key model.ChannelSubscriptionKey, name string) (*chat.ChannelSubscription, error) {

	ctx, l := GetCtx(ctx, "ChatStore", "CreateSubscription")
	l = l.With(zap.Object("channel_subscription", key))

	channelSubscription := chat.ChannelSubscription{
		ChannelID: key.ChannelKey.ID,
		UserID:    key.UserKey.String(),
		Name:      name,
	}

	qry := cs.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "channel_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	}).Create(&channelSubscription)

	l.Info("rows affected", zap.Int64("rows", qry.RowsAffected))

	err := qry.Error

	if err != nil {
		l.Error("could not store channel subscription in database", zap.Error(err))
		return nil, err
	}

	return &channelSubscription, nil
}

func (cs *ChatStore) GetChannel(ctx context.Context, channelKey model.ChannelKey) (*chat.Channel, error) {

	ctx, l := GetCtx(ctx, "ChatStore", "GetChannel")
	l = l.With(zap.String("channelId", channelKey.ID))

	var channel chat.Channel
	err := cs.db.Where("id = ?", channelKey.String()).First(&channel).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		l.Info("channel not found")
		return nil, chat.ErrChannelNotFound
	}

	if err != nil {
		l.Error("could not get channel", zap.Error(err))
		return nil, err
	}

	return &channel, nil

}

func (cs *ChatStore) GetSubscriptionsForUser(ctx context.Context, request *chat.GetSubscriptions) (*chat.ChannelSubscriptions, error) {
	ctx, l := GetCtx(ctx, "ChatStore", "GetSubscriptionsForUser")

	var subscriptions []chat.ChannelSubscription
	err := cs.db.
		Where("user_id = ?", request.UserKey.String()).
		Order("last_message_at desc").
		Offset(request.Skip).
		Limit(request.Take).
		Find(&subscriptions).
		Error

	if err != nil {
		l.Error("could not get subscriptions", zap.Error(err))
		return nil, err
	}
	return chat.NewChannelSubscriptions(subscriptions), nil
}

func (cs *ChatStore) GetSubscriptionsForChannel(ctx context.Context, channelKey model.ChannelKey) ([]chat.ChannelSubscription, error) {
	ctx, _ = GetCtx(ctx, "ChatStore", "GetSubscriptionsForChannel")

	var subscriptions []chat.ChannelSubscription
	err := cs.db.
		Where("channel_id = ?", channelKey.String()).
		Find(&subscriptions).
		Error

	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (cs *ChatStore) GetSubscription(ctx context.Context, request *chat.GetSubscription) (*chat.ChannelSubscription, error) {

	ctx, l := GetCtx(ctx, "ChatStore", "GetSubscription")

	l = l.With(
		zap.String("user_id", request.SubscriptionKey.UserKey.String()),
		zap.String("channel_id", request.SubscriptionKey.ChannelKey.String()))

	l.Debug("getting subscriptions")

	subscription := chat.ChannelSubscription{}

	err := cs.db.First(&subscription, "channel_id = ? and user_id = ?",
		request.SubscriptionKey.ChannelKey.String(),
		request.SubscriptionKey.UserKey.String()).
		Error

	if err != nil {
		l.Error("could not get subscriptions", zap.Error(err))
		return nil, err
	}

	return &subscription, nil

}

func (cs *ChatStore) GetMessage(ctx context.Context, messageKey model.MessageKey) (*chat.Message, error) {

	ctx, l := GetCtx(ctx, "ChatStore", "GetMessage")

	var message Message
	err := cs.db.
		Model(Message{}).
		Where("id = ?", messageKey.String()).
		First(&message).
		Error

	if err != nil {
		l.Error("could not get message", zap.Error(err))
		return nil, err
	}

	returnMessage, err := mapMessage(ctx, &message)
	if err != nil {
		l.Error("could not map message", zap.Error(err))
		return nil, err
	}

	return returnMessage, nil
}

func (cs *ChatStore) GetMessages(ctx context.Context, request *chat.GetMessages) (*chat.GetMessagesResponse, error) {

	ctx, l := GetCtx(ctx, "ChatStore", "GetMessages")

	var messages []Message

	err := cs.db.
		Model(Message{}).
		Where("channel_id = ? AND (visible_to_user IS NULL OR visible_to_user = ?) AND sent_at < ?",
			request.Channel.String(),
			request.UserKey.String(),
			request.Before).
		Order("sent_at desc").
		Limit(request.Take + 1).
		Find(&messages).
		Error

	if err != nil {
		l.Error("could not get messages", zap.Error(err))
		return nil, err
	}

	messageCount := len(messages)
	if messageCount > 0 {
		lastMessageTs := messages[0].SentAt
		err = cs.db.Model(&chat.ChannelSubscription{}).
			Where("channel_id = ? AND user_id = ?",
				request.Channel.String(),
				request.UserKey.String(),
				lastMessageTs).
			Update("last_time_read", lastMessageTs).
			Error
		if err != nil {
			l.Error("could not update subscription", zap.Error(err))
			return nil, err
		}
	}

	if messageCount > request.Take && request.Take > 0 {
		messages = messages[:messageCount-1]
	}

	var mappedMessages []chat.Message
	for _, message := range messages {
		mappedMessage, err := mapMessage(ctx, &message)
		if err != nil {
			l.Error("could not map message", zap.Error(err))
			return nil, err
		}
		mappedMessages = append(mappedMessages, *mappedMessage)
	}

	messageLst := chat.NewMessages(mappedMessages)

	return &chat.GetMessagesResponse{
		Messages: messageLst,
		HasMore:  messageCount > request.Take,
	}, nil
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

func (cs *ChatStore) SaveMessage(ctx context.Context, request *chat.SaveMessageRequest) (*chat.Message, error) {

	ctx, l := GetCtx(ctx, "ChatStore", "SaveMessage")

	l = l.With(
		zap.String("from_user_id", request.FromUser.String()),
		zap.String("channel_id", request.ChannelKey.String()),
		zap.String("text", request.Text))

	var err error
	sentAt := time.Now()

	var subscriptions []chat.ChannelSubscription
	err = cs.db.Model(chat.ChannelSubscription{}).
		Where("channel_id = ?", request.ChannelKey.String()).
		Find(&subscriptions).
		Error

	if err != nil {
		l.Error("could not get subscriptions", zap.Error(err))
		return nil, err
	}

	err = cs.db.Model(&chat.ChannelSubscription{}).
		Where("channel_id = ?", request.ChannelKey.String()).
		Updates(map[string]interface{}{
			"last_message_at":        sentAt,
			"last_message_chars":     utils.FirstChars(request.Text, 30),
			"last_message_user_id":   request.FromUser.String(),
			"last_message_user_name": request.FromUserName,
		}).
		Error

	if err != nil {
		l.Error("could not update subscriptions", zap.Error(err))
		return nil, err
	}

	err = cs.db.Model(&chat.ChannelSubscription{}).
		Where("channel_id = ? and user_id = ?",
			request.ChannelKey.String(),
			request.FromUser.String()).
		Updates(map[string]interface{}{
			"last_time_read": sentAt,
		}).
		Error

	if err != nil {
		l.Error("could not update sender channel subscription")
		return nil, err
	}

	var visibleToUser *string = nil
	if request.VisibleToUser != nil {

		visibleToUserStr := request.VisibleToUser.String()
		visibleToUser = &visibleToUserStr
	}

	messageKey := model.NewMessageKey(uuid.NewV4())

	blocksJson := []byte("[]")
	if request.Blocks != nil && len(request.Blocks) > 0 {
		blocksJson, err = json.Marshal(request.Blocks)
		if err != nil {
			l.Error("could not unmarshal message blocks", zap.Error(err))
			return nil, err
		}
	}

	attachmentsJson := []byte("[]")
	if request.Attachments != nil && len(request.Attachments) > 0 {
		attachmentsJson, err = json.Marshal(request.Attachments)
		if err != nil {
			l.Error("could not unmarshal message attachments", zap.Error(err))
			return nil, err
		}
	}

	message := Message{
		ID:             messageKey.GetUUID(),
		ChannelID:      request.ChannelKey.String(),
		MessageType:    chat.NormalMessage,
		MessageSubType: chat.UserMessage,
		SentById:       request.FromUser.String(),
		SentByUsername: request.FromUserName,
		SentAt:         sentAt,
		Text:           request.Text,
		Blocks:         string(blocksJson),
		Attachments:    string(attachmentsJson),
		VisibleToUser:  visibleToUser,
	}

	err = cs.db.Create(message).Error
	if err != nil {
		l.Error("could not save message", zap.Error(err))
		return nil, err
	}

	returnMessage, err := mapMessage(ctx, &message)
	if err != nil {
		l.Error("could not map message", zap.Error(err))
	}

	return returnMessage, nil

}

func mapMessage(ctx context.Context, message *Message) (*chat.Message, error) {

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

func (cs *ChatStore) GetTopic(key model.ChannelKey) (*chat.Channel, error) {
	var topic chat.Channel
	err := cs.db.Where("id = ?", key.ID).First(&topic).Error
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

func (cs *ChatStore) getChannel(channelKey *model.ChannelKey) (*chat.Channel, error) {
	var channel chat.Channel
	err := cs.db.
		Where("id = ?", channelKey.String()).
		First(&channel).
		Error
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (cs *ChatStore) getConversationId(participants *model.UserKeys) (*model.ChannelKey, error) {
	var shortUids []string

	for _, participant := range participants.Items {
		sid, err := utils.ShortUuidFromStr(participant.String())
		if err != nil {
			return nil, err
		}
		shortUids = append(shortUids, sid)
	}

	sort.Strings(shortUids)
	channelId := strings.Join(shortUids, "-")
	channelKey := model.NewConversationKey(channelId)
	return &channelKey, nil
}

func (cs *ChatStore) getGroupChannelKey(group *model.GroupKey) (*model.ChannelKey, error) {
	shortUuid := utils.ShortUuid(group.ID)
	channelKey := model.NewConversationKey(shortUuid)
	return &channelKey, nil
}

func (cs *ChatStore) getLogger(ctx context.Context) *zap.Logger {
	return logging.WithContext(ctx).Named("ChatStore")
}
