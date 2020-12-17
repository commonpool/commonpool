package store

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/keys"
)

func (cs *ChatStore) CreateSubscription(ctx context.Context, key keys.ChannelSubscriptionKey, name string) (*chat.ChannelSubscription, error) {

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
