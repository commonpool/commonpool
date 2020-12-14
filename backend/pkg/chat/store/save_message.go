package store

import (
	"context"
	"encoding/json"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/utils"
	"time"
)

func (cs *ChatStore) SaveMessage(ctx context.Context, request *chat.Message) error {

	var err error
	sentAt := time.Now()

	var subscriptions []ChannelSubscription
	err = cs.db.Model(ChannelSubscription{}).
		Where("channel_id = ?", request.ChannelKey.String()).
		Find(&subscriptions).
		Error

	if err != nil {
		return err
	}

	err = cs.db.Model(&ChannelSubscription{}).
		Where("channel_id = ?", request.ChannelKey.String()).
		Updates(map[string]interface{}{
			"last_message_at":        sentAt,
			"last_message_chars":     utils.FirstChars(request.Text, 30),
			"last_message_user_id":   request.SentBy.UserKey.String(),
			"last_message_user_name": request.SentBy.Username,
		}).
		Error

	if err != nil {
		return err
	}

	err = cs.db.Model(&ChannelSubscription{}).
		Where("channel_id = ? and user_id = ?",
			request.ChannelKey.String(),
			request.SentBy.UserKey.String()).
		Updates(map[string]interface{}{
			"last_time_read": sentAt,
		}).
		Error

	if err != nil {
		return err
	}

	var visibleToUser *string = nil
	if request.VisibleToUser != nil {

		visibleToUserStr := request.VisibleToUser.String()
		visibleToUser = &visibleToUserStr
	}

	blocksJson := []byte("[]")
	if request.Blocks != nil && len(request.Blocks) > 0 {
		blocksJson, err = json.Marshal(request.Blocks)
		if err != nil {
			return err
		}
	}

	attachmentsJson := []byte("[]")
	if request.Attachments != nil && len(request.Attachments) > 0 {
		attachmentsJson, err = json.Marshal(request.Attachments)
		if err != nil {
			return err
		}
	}

	message := Message{
		ID:             request.Key.GetUUID(),
		ChannelID:      request.ChannelKey.String(),
		MessageType:    chat.NormalMessage,
		MessageSubType: chat.UserMessage,
		SentById:       request.SentBy.UserKey.String(),
		SentByUsername: request.SentBy.Username,
		SentAt:         sentAt,
		Text:           request.Text,
		Blocks:         string(blocksJson),
		Attachments:    string(attachmentsJson),
		VisibleToUser:  visibleToUser,
	}

	err = cs.db.Create(message).Error
	if err != nil {
		return err
	}

	return nil

}
