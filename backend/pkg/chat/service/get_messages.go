package service

import (
	"context"
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/chat"
	"time"
)

func (c ChatService) GetMessages(ctx context.Context, channel chat.ChannelKey, before time.Time, take int) (*chat.GetMessagesResponse, error) {

	loggedInUser, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}

	return c.chatStore.GetMessages(ctx, &chat.GetMessages{
		Take:    take,
		Before:  before,
		Channel: channel,
		UserKey: loggedInUser.GetUserKey(),
	})
}
