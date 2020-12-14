package service

import (
	"context"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/chat"
)

// SubscribeToChannel will subscribe a user to a given channel
func (c ChatService) SubscribeToChannel(ctx context.Context, channelSubscriptionKey model.ChannelSubscriptionKey, name string) (*chat.ChannelSubscription, error) {

	channelSubscription, err := c.chatStore.CreateSubscription(ctx, channelSubscriptionKey, name)
	if err != nil {
		return nil, err
	}

	amqpChannel, err := c.amqpClient.GetChannel()
	if err != nil {
		return nil, err
	}
	defer amqpChannel.Close()

	userExchangeName, err := c.CreateUserExchange(ctx, channelSubscriptionKey.UserKey)
	if err != nil {
		return nil, err
	}

	headers := c.getChannelBindingHeaders(channelSubscriptionKey)
	err = amqpChannel.ExchangeBind(ctx, userExchangeName, "", amqp.WebsocketMessagesExchange, false, headers)
	if err != nil {
		return nil, err
	}

	return channelSubscription, nil

}
