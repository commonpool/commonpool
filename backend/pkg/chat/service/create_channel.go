package service

import (
	"context"
	"github.com/commonpool/backend/pkg/chat/chatmodel"
)

func (c ChatService) CreateChannel(ctx context.Context, channelKey chatmodel.ChannelKey, channelType chatmodel.ChannelType) (*chatmodel.Channel, error) {
	channel := &chatmodel.Channel{
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
