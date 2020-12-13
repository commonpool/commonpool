package store

import (
	"context"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/store"
	"go.uber.org/zap"
)

func (cs *ChatStore) GetSubscription(ctx context.Context, request *chat.GetSubscription) (*chat.ChannelSubscription, error) {

	ctx, l := store.GetCtx(ctx, "ChatStore", "GetSubscription")

	l = l.With(
		zap.String("user_id", request.SubscriptionKey.UserKey.String()),
		zap.String("channel_id", request.SubscriptionKey.ChannelKey.String()))

	l.Debug("getting subscriptions")

	subscription := chat.ChannelSubscription{}

	err := cs.db.First(&subscription, "channel_id = ? and user_id = ?",
		request.SubscriptionKey.ChannelKey.String(),
		request.SubscriptionKey.UserKey.String()).
		Error

	if err != nil {
		l.Error("could not get subscriptions", zap.Error(err))
		return nil, err
	}

	return &subscription, nil

}
