package store

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/chat/model"
)

func (cs *ChatStore) GetSubscription(ctx context.Context, request *chat.GetSubscription) (*model.ChannelSubscription, error) {

	subscription := model.ChannelSubscription{}

	err := cs.db.First(&subscription, "channel_id = ? and user_id = ?",
		request.SubscriptionKey.ChannelKey.String(),
		request.SubscriptionKey.UserKey.String()).
		Error

	if err != nil {
		return nil, err
	}

	return &subscription, nil

}
