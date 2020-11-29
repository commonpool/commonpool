package service

import (
	"context"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
	"go.uber.org/zap"
)

func (c ChatService) CreateChannel(ctx context.Context, channelKey model.ChannelKey, channelType chat.ChannelType) (*chat.Channel, error) {

	ctx, l := GetCtx(ctx, "ChatService", "CreateChannel")

	channel := &chat.Channel{
		ID:    channelKey.ID,
		Title: "",
		Type:  channelType,
	}

	err := c.cs.CreateChannel(ctx, channel)
	if err != nil {
		l.Error("could not create channel", zap.Error(err))
		return nil, err
	}

	channel, err = c.cs.GetChannel(ctx, channelKey)
	if err != nil {
		l.Error("could not get channel", zap.Error(err))
		return nil, err
	}

	return channel, nil
}
