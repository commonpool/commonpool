package service

import (
	"context"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/model"
)

func (c ChatService) UnsubscribeFromChannel(ctx context.Context, channelSubscriptionKey model.ChannelSubscriptionKey) error {

	err := c.chatStore.DeleteSubscription(ctx, channelSubscriptionKey)
	if err != nil {
		return err
	}

	amqpChannel, err := c.amqpClient.GetChannel()
	if err != nil {
		return err
	}
	defer amqpChannel.Close()

	userExchangeName := channelSubscriptionKey.UserKey.GetExchangeName()
	headers := c.getChannelBindingHeaders(channelSubscriptionKey)
	err = amqpChannel.ExchangeUnbind(ctx, userExchangeName, "", amqp.WebsocketMessagesExchange, false, headers)
	if err != nil {
		return err
	}

	return nil

}
