package service

import (
	"context"
	"github.com/commonpool/backend/model"
	"go.uber.org/zap"
)

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
