package service

import (
	"context"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
)

func (c ChatService) CreateChannel(ctx context.Context, channelKey model.ChannelKey, channelType chat.ChannelType) (*chat.Channel, error) {

	channel := &chat.Channel{
		ID:    channelKey.ID,
		Title: "",
		Type:  channelType,
	}

	err := c.cs.CreateChannel(ctx, channel)
	if err != nil {
		return nil, err
	}

	channel, err = c.cs.GetChannel(ctx, channelKey)
	if err != nil {
		return nil, err
	}

	return channel, nil
}
