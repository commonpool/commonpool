package store

import (
	"context"
	"github.com/commonpool/backend/chat"
)

func (cs *ChatStore) CreateChannel(ctx context.Context, channel *chat.Channel) error {
	return cs.db.WithContext(ctx).Create(channel).Error
}
