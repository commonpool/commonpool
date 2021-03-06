package handler

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/chat/service"
	"time"
)

type Subscription struct {
	ChannelID           string           `json:"channelId"`
	UserID              string           `json:"userId"`
	HasUnreadMessages   bool             `json:"hasUnreadMessages"`
	CreatedAt           time.Time        `json:"createdAt"`
	UpdatedAt           time.Time        `json:"updatedAt"`
	LastMessageAt       time.Time        `json:"lastMessageAt"`
	LastTimeRead        time.Time        `json:"lastTimeRead"`
	LastMessageChars    string           `json:"lastMessageChars"`
	LastMessageUserId   string           `json:"lastMessageUserId"`
	LastMessageUserName string           `json:"lastMessageUsername"`
	Name                string           `json:"name"`
	Type                chat.ChannelType `json:"type"`
}

func MapSubscription(channel *chat.Channel, subscription *chat.ChannelSubscription) *Subscription {
	return &Subscription{
		ChannelID:           channel.Key.String(),
		UserID:              subscription.UserKey.String(),
		HasUnreadMessages:   subscription.LastMessageAt.After(subscription.LastTimeRead),
		CreatedAt:           subscription.CreatedAt,
		UpdatedAt:           subscription.UpdatedAt,
		LastMessageAt:       subscription.LastMessageAt,
		LastTimeRead:        subscription.LastTimeRead,
		LastMessageChars:    subscription.LastMessageChars,
		LastMessageUserId:   subscription.LastMessageUserKey.String(),
		LastMessageUserName: subscription.LastMessageUserName,
		Name:                subscription.Name,
		Type:                channel.Type,
	}
}

func MapSubscriptions(ctx context.Context, chatService service.Service, subscriptions *chat.ChannelSubscriptions) ([]Subscription, error) {
	var items []Subscription
	for _, subscription := range subscriptions.Items {
		channel, err := chatService.GetChannel(ctx, subscription.GetChannelKey())
		if err != nil {
			return nil, err
		}
		items = append(items, *MapSubscription(channel, subscription))
	}
	if items == nil {
		items = []Subscription{}
	}
	return items, nil
}
