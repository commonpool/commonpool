package service

import (
	"context"
	chatmodel "github.com/commonpool/backend/pkg/chat/model"
	"github.com/commonpool/backend/pkg/mq"
)

// SubscribeToChannel will subscribe a user to a given channel
func (c ChatService) SubscribeToChannel(ctx context.Context, channelSubscriptionKey chatmodel.ChannelSubscriptionKey, name string) (*chatmodel.ChannelSubscription, error) {

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
	err = amqpChannel.ExchangeBind(ctx, userExchangeName, "", mq.WebsocketMessagesExchange, false, headers)
	if err != nil {
		return nil, err
	}

	return channelSubscription, nil

}
