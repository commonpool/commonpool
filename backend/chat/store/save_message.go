package store

import (
	"context"
	"encoding/json"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/store"
	"github.com/commonpool/backend/utils"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"time"
)

func (cs *ChatStore) SaveMessage(ctx context.Context, request *chat.SaveMessageRequest) (*chat.Message, error) {

	ctx, l := store.GetCtx(ctx, "ChatStore", "SaveMessage")

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

	message := store.Message{
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
