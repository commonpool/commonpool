package service

import (
	"context"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
	"time"
)

func (c ChatService) GetMessages(ctx context.Context, userKey model.UserKey, channel model.ChannelKey, before time.Time, take int) (*chat.GetMessagesResponse, error) {
	return c.cs.GetMessages(ctx, &chat.GetMessages{
		Take:    take,
		Before:  before,
		Channel: channel,
		UserKey: userKey,
	})
}
