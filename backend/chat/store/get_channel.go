package store

import (
	"context"
	"errors"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/store"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (cs *ChatStore) GetChannel(ctx context.Context, channelKey model.ChannelKey) (*chat.Channel, error) {

	ctx, l := store.GetCtx(ctx, "ChatStore", "GetChannel")
	l = l.With(zap.String("channelId", channelKey.ID))

	var channel chat.Channel
	err := cs.db.Where("id = ?", channelKey.String()).First(&channel).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		l.Info("channel not found")
		return nil, chat.ErrChannelNotFound
	}

	if err != nil {
		l.Error("could not get channel", zap.Error(err))
		return nil, err
	}

	return &channel, nil

}
