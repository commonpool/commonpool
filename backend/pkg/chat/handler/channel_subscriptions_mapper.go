package handler

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/web"
)

func MapChannelSubscriptions(ctx context.Context, chatService chat.Service, subscriptions *chat.ChannelSubscriptions) ([]web.Subscription, error) {

	var items []web.Subscription
	for _, subscription := range subscriptions.Items {
		channel, err := chatService.GetChannel(ctx, subscription.GetChannelKey())
		if err != nil {
			return nil, err
		}
		items = append(items, *web.MapSubscription(channel, &subscription))
	}

	if items == nil {
		items = []web.Subscription{}
	}

	return items, nil

}
