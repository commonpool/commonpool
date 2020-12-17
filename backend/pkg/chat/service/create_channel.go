package service

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/keys"
)

func (c ChatService) CreateChannel(ctx context.Context, channelKey keys.ChannelKey, channelType chat.ChannelType) (*chat.Channel, error) {
	channel := &chat.Channel{
		Key:   channelKey,
		Title: "",
		Type:  channelType,
	}

	err := c.chatStore.CreateChannel(ctx, channel)
	if err != nil {
		return nil, err
	}

	channel, err = c.chatStore.GetChannel(ctx, channelKey)
	if err != nil {
		return nil, err
	}

	return channel, nil
}
