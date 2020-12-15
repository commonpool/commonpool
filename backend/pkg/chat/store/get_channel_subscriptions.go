package store

import (
	"context"
	"github.com/commonpool/backend/pkg/chat/model"
)

func (cs *ChatStore) GetSubscriptionsForChannel(ctx context.Context, channelKey model.ChannelKey) ([]*model.ChannelSubscription, error) {

	var subscriptions []ChannelSubscription
	err := cs.db.
		Where("channel_id = ?", channelKey.String()).
		Find(&subscriptions).
		Error

	if err != nil {
		return nil, err
	}

	var result []*model.ChannelSubscription
	for _, subscription := range subscriptions {
		mappedSubscription := subscription.Map()
		result = append(result, mappedSubscription)
	}

	return result, nil
}