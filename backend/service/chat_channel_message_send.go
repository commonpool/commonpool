package service

import (
	"context"
	"encoding/json"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
	"go.uber.org/zap"
)

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

	return nil, nil
}
