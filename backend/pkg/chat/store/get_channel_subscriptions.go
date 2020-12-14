package store

import (
	"context"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/chat"
)

func (cs *ChatStore) GetSubscriptionsForChannel(ctx context.Context, channelKey model.ChannelKey) ([]chat.ChannelSubscription, error) {

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
