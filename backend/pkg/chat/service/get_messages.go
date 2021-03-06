package service

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/chat/store"
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

func (c ChatService) GetMessages(ctx context.Context, channel keys.ChannelKey, before time.Time, take int) (*store.GetMessagesResponse, error) {

	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}

	return c.chatStore.GetMessages(ctx, &store.GetMessages{
		Take:    take,
		Before:  before,
		Channel: channel,
		UserKey: loggedInUser.GetUserKey(),
	})
}
