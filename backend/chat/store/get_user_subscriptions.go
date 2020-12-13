package store

import (
	"context"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/store"
	"go.uber.org/zap"
)

func (cs *ChatStore) GetSubscriptionsForUser(ctx context.Context, request *chat.GetSubscriptions) (*chat.ChannelSubscriptions, error) {
	ctx, l := store.GetCtx(ctx, "ChatStore", "GetSubscriptionsForUser")

	var subscriptions []chat.ChannelSubscription
	err := cs.db.
		Where("user_id = ?", request.UserKey.String()).
		Order("last_message_at desc").
		Offset(request.Skip).
		Limit(request.Take).
		Find(&subscriptions).
		Error

	if err != nil {
		l.Error("could not get subscriptions", zap.Error(err))
		return nil, err
	}
	return chat.NewChannelSubscriptions(subscriptions), nil
}
