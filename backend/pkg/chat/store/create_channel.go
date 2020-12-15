package store

import (
	"context"
	"github.com/commonpool/backend/pkg/chat/model"
)

func (cs *ChatStore) CreateChannel(ctx context.Context, channel *model.Channel) error {
	dbChannel := MapChannel(channel)
	return cs.db.WithContext(ctx).Create(dbChannel).Error
}
