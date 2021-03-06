package service

import (
	"context"
	"encoding/json"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/mq"
	uuid "github.com/satori/go.uuid"
	"time"
)

func (c ChatService) SendConversationMessage(ctx context.Context, request *SendConversationMessage) (*SendConversationMessageResponse, error) {

	createdChannel, err := c.getOrCreateConversationChannel(ctx, request.ToUserKeys)
	if err != nil {
		return nil, err
	}

	channelKey := createdChannel.Channel.GetKey()
	message := &chat.Message{
		Key:            keys.NewMessageKey(uuid.NewV4()),
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
		Blocks:        nil,
		Attachments:   nil,
		VisibleToUser: nil,
	}
	err = c.chatStore.SaveMessage(ctx, message)

	if err != nil {
		return nil, err
	}

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

	amqpChannel, err := c.amqpClient.GetChannel()
	if err != nil {
		return nil, err
	}
	defer amqpChannel.Close()

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

	return &SendConversationMessageResponse{
		Message: message,
	}, nil
}
