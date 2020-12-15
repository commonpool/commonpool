package handler

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/chat/handler/model"
	model2 "github.com/commonpool/backend/pkg/chat/model"
)

func MapChannelSubscriptions(ctx context.Context, chatService chat.Service, subscriptions *model2.ChannelSubscriptions) ([]model.Subscription, error) {

	var items []model.Subscription
	for _, subscription := range subscriptions.Items {
		channel, err := chatService.GetChannel(ctx, subscription.GetChannelKey())
		if err != nil {
			return nil, err
		}
		items = append(items, *model.MapSubscription(channel, &subscription))
	}

	if items == nil {
		items = []model.Subscription{}
	}

	return items, nil

}
