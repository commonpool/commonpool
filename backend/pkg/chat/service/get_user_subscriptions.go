package service

import (
	"context"
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/chat/chatmodel"
)

func (c ChatService) GetSubscriptionsForUser(ctx context.Context, take int, skip int) (*chatmodel.ChannelSubscriptions, error) {

	loggedInUser, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return nil, err
	}
	loggedInUserKey := loggedInUser.GetUserKey()

	subs, err := c.chatStore.GetSubscriptionsForUser(ctx, &chat.GetSubscriptions{
		Take:    take,
		Skip:    skip,
		UserKey: loggedInUserKey,
	})

	if err != nil {
		return nil, err
	}

	return chatmodel.NewChannelSubscriptions(subs.Items), nil
}
