package service

import (
	"context"
	"encoding/json"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/service"
	"go.uber.org/zap"
)

func (c ChatService) SendGroupMessage(ctx context.Context, request *chat.SendGroupMessage) (*chat.SendGroupMessageResponse, error) {

	ctx, l := service.GetCtx(ctx, "ChatService", "SendConversationMessage")

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
		ID:        sentMessage.Key.String(),
		Timestamp: sentMessage.SentAt.String(),
		Text:      sentMessage.Text,
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

	// todo
	return &chat.SendGroupMessageResponse{}, nil

}
