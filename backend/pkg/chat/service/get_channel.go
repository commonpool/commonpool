package service

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/keys"
)

func (c ChatService) GetChannel(ctx context.Context, channelKey keys.ChannelKey) (*chat.Channel, error) {
	return c.chatStore.GetChannel(ctx, channelKey)
}
