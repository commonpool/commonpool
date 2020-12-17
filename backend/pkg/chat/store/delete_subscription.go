package store

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/keys"
)

func (cs *ChatStore) DeleteSubscription(ctx context.Context, key keys.ChannelSubscriptionKey) error {

	err := cs.db.Delete(chat.ChannelSubscription{},
		"user_id = ? and channel_id = ?",
		key.UserKey.String(),
		key.ChannelKey.String()).
		Error

	if err != nil {
		return err
	}

	return nil
}
