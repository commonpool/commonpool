package store

import (
	"context"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/store"
)

func (cs *ChatStore) GetSubscriptionsForChannel(ctx context.Context, channelKey model.ChannelKey) ([]chat.ChannelSubscription, error) {
	ctx, _ = store.GetCtx(ctx, "ChatStore", "GetSubscriptionsForChannel")

	var subscriptions []chat.ChannelSubscription
	err := cs.db.
		Where("channel_id = ?", channelKey.String()).
		Find(&subscriptions).
		Error

	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}
