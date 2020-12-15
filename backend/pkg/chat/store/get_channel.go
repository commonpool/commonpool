package store

import (
	"context"
	"errors"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/chat/model"
	"gorm.io/gorm"
)

func (cs *ChatStore) GetChannel(ctx context.Context, channelKey model.ChannelKey) (*model.Channel, error) {

	var channel Channel
	err := cs.db.Where("id = ?", channelKey.String()).First(&channel).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, chat.ErrChannelNotFound
	}

	if err != nil {
		return nil, err
	}

	return channel.Map(), nil

}