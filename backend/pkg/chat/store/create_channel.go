package store

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
)

func (cs *ChatStore) CreateChannel(ctx context.Context, channel *chat.Channel) error {
	dbChannel := MapChannel(channel)
	return cs.db.WithContext(ctx).Create(dbChannel).Error
}
