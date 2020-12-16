package store

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/chat/chatmodel"
)

func (cs *ChatStore) GetSubscriptionsForUser(ctx context.Context, request *chat.GetSubscriptions) (*chatmodel.ChannelSubscriptions, error) {

	var subscriptions []chatmodel.ChannelSubscription
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
	return chatmodel.NewChannelSubscriptions(subscriptions), nil
}
