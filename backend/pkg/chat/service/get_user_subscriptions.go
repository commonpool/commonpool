package service

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/chat/store"
)

func (c ChatService) GetSubscriptionsForUser(ctx context.Context, take int, skip int) (*chat.ChannelSubscriptions, error) {

	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	loggedInUserKey := loggedInUser.GetUserKey()

	subs, err := c.chatStore.GetSubscriptionsForUser(ctx, &store.GetSubscriptions{
		Take:    take,
		Skip:    skip,
		UserKey: loggedInUserKey,
	})

	if err != nil {
		return nil, err
	}

	return chat.NewChannelSubscriptions(subs.Items), nil
}
