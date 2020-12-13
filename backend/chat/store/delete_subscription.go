package store

import (
	"context"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/store"
	"go.uber.org/zap"
)

func (cs *ChatStore) DeleteSubscription(ctx context.Context, key model.ChannelSubscriptionKey) error {

	ctx, l := store.GetCtx(ctx, "ChatStore", "CreateSubscription")
	l = l.With(zap.Object("channel_subscription", key))

	err := cs.db.Delete(chat.ChannelSubscription{},
		"user_id = ? and channel_id = ?",
		key.UserKey.String(),
		key.ChannelKey.String()).
		Error

	if err != nil {
		l.Error("could not delete channel subscription", zap.Error(err))
		return err
	}

	return nil
}
