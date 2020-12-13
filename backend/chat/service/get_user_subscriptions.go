package service

import (
	"context"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
)

func (c ChatService) GetUserSubscriptions(ctx context.Context, userKey model.UserKey, take int, skip int) (*chat.ChannelSubscriptions, error) {

	subs, err := c.cs.GetSubscriptionsForUser(ctx, &chat.GetSubscriptions{
		Take:    take,
		Skip:    skip,
		UserKey: userKey,
	})

	if err != nil {
		return nil, err
	}

	return chat.NewChannelSubscriptions(subs.Items), nil
}
