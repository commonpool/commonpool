package service

import (
	"context"
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/chat/model"
	"time"
)

func (c ChatService) GetMessages(ctx context.Context, channel model.ChannelKey, before time.Time, take int) (*chat.GetMessagesResponse, error) {

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
