package service

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/mq"
)

// SubscribeToChannel will subscribe a user to a given channel
func (c ChatService) SubscribeToChannel(ctx context.Context, channelSubscriptionKey keys.ChannelSubscriptionKey, name string) (*chat.ChannelSubscription, error) {

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
