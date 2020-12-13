package service

import (
	"context"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
)

func (c ChatService) GetChannel(ctx context.Context, channelKey model.ChannelKey) (*chat.Channel, error) {
	return c.cs.GetChannel(ctx, channelKey)
}
