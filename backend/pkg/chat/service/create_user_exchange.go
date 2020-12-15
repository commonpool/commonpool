package service

import (
	"context"
	usermodel "github.com/commonpool/backend/pkg/user/model"
)

// CreateUserExchange will create the AMQP exchange to receive user messages
// This exchange will be bound to queues representing different Websocket clients
// for the same user (if a user is using multiple devices to connect, he
// will get Websocket notifications on all devices)
func (c ChatService) CreateUserExchange(ctx context.Context, userKey usermodel.UserKey) (string, error) {

	amqpChannel, err := c.amqpClient.GetChannel()
	if err != nil {
		return "", err
	}
	defer amqpChannel.Close()

	exchangeName := c.GetUserExchangeName(ctx, userKey)

	err = amqpChannel.ExchangeDeclare(ctx, exchangeName, "fanout", true, false, false, false, nil)
	if err != nil {
		return "", err
	}

	return exchangeName, nil
}
