package store

import (
	"context"
	chatmodel "github.com/commonpool/backend/pkg/chat/chatmodel"
)

func (cs *ChatStore) DeleteSubscription(ctx context.Context, key chatmodel.ChannelSubscriptionKey) error {

	err := cs.db.Delete(chatmodel.ChannelSubscription{},
		"user_id = ? and channel_id = ?",
		key.UserKey.String(),
		key.ChannelKey.String()).
		Error

	if err != nil {
		return err
	}

	return nil
}
