package store

import (
	"context"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/store"
	"go.uber.org/zap"
)

func (cs *ChatStore) CreateSubscription(ctx context.Context, key model.ChannelSubscriptionKey, name string) (*chat.ChannelSubscription, error) {

	ctx, l := store.GetCtx(ctx, "ChatStore", "CreateSubscription")
	l = l.With(zap.Object("channel_subscription", key))

	channelSubscription := chat.ChannelSubscription{
		ChannelID: key.ChannelKey.ID,
		UserID:    key.UserKey.String(),
		Name:      name,
	}

	err := cs.db.Create(&channelSubscription).Error

	if err != nil {
		l.Error("could not store channel subscription in database", zap.Error(err))
		return nil, err
	}

	return &channelSubscription, nil
}
