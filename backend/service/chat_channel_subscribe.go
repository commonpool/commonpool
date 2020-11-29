package service

import (
	"context"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
	"go.uber.org/zap"
)

// SubscribeToChannel will subscribe a user to a given channel
func (c ChatService) SubscribeToChannel(ctx context.Context, channelSubscriptionKey model.ChannelSubscriptionKey, name string) (*chat.ChannelSubscription, error) {

	ctx, l := GetCtx(ctx, "ChatService", "SubscribeToChannel")
	l = l.With(zap.Object("channel_subscription", channelSubscriptionKey))

	l.Debug("subscribing to channel")

	channelSubscription, err := c.cs.CreateSubscription(ctx, channelSubscriptionKey, name)
	if err != nil {
		l.Error("could not create channel subscription", zap.Error(err))
		return nil, err
	}

	amqpChannel, err := c.mq.GetChannel()
	if err != nil {
		l.Error("could not get amqp channel", zap.Error(err))
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
		l.Error("could not bind user exchange", zap.Error(err))
		return nil, err
	}

	return channelSubscription, nil

}
