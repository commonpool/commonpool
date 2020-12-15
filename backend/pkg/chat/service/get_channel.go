package service

import (
	"context"
	"github.com/commonpool/backend/pkg/chat/model"
)

func (c ChatService) GetChannel(ctx context.Context, channelKey model.ChannelKey) (*model.Channel, error) {
	return c.chatStore.GetChannel(ctx, channelKey)
}
