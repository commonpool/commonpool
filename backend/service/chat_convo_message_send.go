package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/chat"
	"go.uber.org/zap"
)

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

	err = amqpChannel.Publish(ctx, amqp.MessagesExchange, "", false, false, amqp.Publishing{
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
