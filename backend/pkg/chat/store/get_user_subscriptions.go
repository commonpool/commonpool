package store

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
)

func (cs *ChatStore) GetSubscriptionsForUser(ctx context.Context, request *chat.GetSubscriptions) (*chat.ChannelSubscriptions, error) {

	var subscriptions []ChannelSubscription
	err := cs.db.
		Where("user_id = ?", request.UserKey.String()).
		Order("last_message_at desc").
		Offset(request.Skip).
		Limit(request.Take).
		Find(&subscriptions).
		Error

	if err != nil {
		return nil, err
	}

	var result []*chat.ChannelSubscription
	for _, subscription := range subscriptions {
		result = append(result, subscription.Map())
	}

	return chat.NewChannelSubscriptions(result), nil
}
