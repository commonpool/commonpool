package service

import (
	"context"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/chat"
)

func (c ChatService) GetChannel(ctx context.Context, channelKey model.ChannelKey) (*chat.Channel, error) {
	return c.chatStore.GetChannel(ctx, channelKey)
}
