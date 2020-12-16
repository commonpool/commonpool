package store

import (
	"context"
	"github.com/commonpool/backend/pkg/chat/chatmodel"
)

func (cs *ChatStore) CreateChannel(ctx context.Context, channel *chatmodel.Channel) error {
	dbChannel := MapChannel(channel)
	return cs.db.WithContext(ctx).Create(dbChannel).Error
}
