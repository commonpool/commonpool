package service

import (
	"context"
	"encoding/json"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/mq"
	uuid "github.com/satori/go.uuid"
	"time"
)

func (c ChatService) SendGroupMessage(ctx context.Context, request *chat.SendGroupMessage) (*chat.SendGroupMessageResponse, error) {

	channelKey := chat.GetChannelKeyForGroup(request.GroupKey)

	message := &chat.Message{
		Key:            chat.NewMessageKey(uuid.NewV4()),
		ChannelKey:     channelKey,
		MessageType:    chat.NormalMessage,
		MessageSubType: chat.UserMessage,
		SentBy: chat.MessageSender{
			Type:     chat.UserMessageSender,
			UserKey:  request.FromUserKey,
			Username: request.FromUserName,
		},
		SentAt:        time.Now(),
		Text:          request.Text,
		Blocks:        request.Blocks,
		Attachments:   request.Attachments,
		VisibleToUser: request.OnlyVisibleToUserKey,
	}

	err := c.chatStore.SaveMessage(ctx, message)

	if err != nil {
		return nil, err
	}

	amqpChannel, err := c.amqpClient.GetChannel()
	if err != nil {
		return nil, err
	}
	defer amqpChannel.Close()

	evt := mq.Event{
		Type:      "message",
		SubType:   "user",
		Channel:   channelKey.String(),
		User:      request.FromUserKey.String(),
		ID:        message.Key.String(),
		Timestamp: message.SentAt.String(),
		Text:      message.Text,
	}

	js, err := json.Marshal(evt)
	if err != nil {
		return nil, err
	}

	err = amqpChannel.Publish(ctx, mq.MessagesExchange, "", false, false, mq.Message{
		Headers: map[string]interface{}{
			"channel_id": channelKey.String(),
			"event_type": "chat.message",
		},
		ContentType: "application/json",
		Body:        js,
	})

	if err != nil {
		return nil, err
	}

	// todo
	return &chat.SendGroupMessageResponse{}, nil

}
