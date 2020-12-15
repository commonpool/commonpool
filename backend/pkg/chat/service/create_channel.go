package service

import (
	"context"
	"github.com/commonpool/backend/pkg/chat/model"
)

func (c ChatService) CreateChannel(ctx context.Context, channelKey model.ChannelKey, channelType model.ChannelType) (*model.Channel, error) {
	channel := &model.Channel{
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
