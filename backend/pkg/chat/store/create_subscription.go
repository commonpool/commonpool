package store

import (
	"context"
	chatmodel "github.com/commonpool/backend/pkg/chat/model"
)

func (cs *ChatStore) CreateSubscription(ctx context.Context, key chatmodel.ChannelSubscriptionKey, name string) (*chatmodel.ChannelSubscription, error) {

	channelSubscription := ChannelSubscription{
		ChannelID: key.ChannelKey.String(),
		UserID:    key.UserKey.String(),
		Name:      name,
	}

	err := cs.db.Create(&channelSubscription).Error

	if err != nil {
		return nil, err
	}

	return channelSubscription.Map(), nil
}
