package service

import (
	"context"
	"github.com/commonpool/backend/pkg/chat/chatmodel"
)

func (c ChatService) GetChannel(ctx context.Context, channelKey chatmodel.ChannelKey) (*chatmodel.Channel, error) {
	return c.chatStore.GetChannel(ctx, channelKey)
}
