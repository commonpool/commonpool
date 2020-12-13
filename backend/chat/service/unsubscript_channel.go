package service

import (
	"context"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/service"
	"go.uber.org/zap"
)

func (c ChatService) UnsubscribeFromChannel(ctx context.Context, channelSubscriptionKey model.ChannelSubscriptionKey) error {

	ctx, l := service.GetCtx(ctx, "ChatService", "UnsubscribeFromChannel")
	l = l.With(zap.Object("channel_subscription", channelSubscriptionKey))

	l.Debug("unsubscribing from channel")

	err := c.cs.DeleteSubscription(ctx, channelSubscriptionKey)
	if err != nil {
		l.Error("could not delete channel subscription", zap.Error(err))
		return err
	}

	amqpChannel, err := c.mq.GetChannel()
	if err != nil {
		l.Error("could not get amqp channel", zap.Error(err))
		return err
	}
	defer amqpChannel.Close()

	userExchangeName := channelSubscriptionKey.UserKey.GetExchangeName()
	headers := c.getChannelBindingHeaders(channelSubscriptionKey)
	err = amqpChannel.ExchangeUnbind(ctx, userExchangeName, "", amqp.WebsocketMessagesExchange, false, headers)
	if err != nil {
		l.Error("could not bind user exchange", zap.Error(err))
		return err
	}

	return nil

}
