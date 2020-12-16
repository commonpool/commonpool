package service

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/mq"
)

func (c ChatService) UnsubscribeFromChannel(ctx context.Context, channelSubscriptionKey chat.ChannelSubscriptionKey) error {

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
	err = amqpChannel.ExchangeUnbind(ctx, userExchangeName, "", mq.WebsocketMessagesExchange, false, headers)
	if err != nil {
		return err
	}

	return nil

}
