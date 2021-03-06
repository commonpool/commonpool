package store

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type Store interface {
	GetSubscriptionsForUser(ctx context.Context, request *GetSubscriptions) (*chat.ChannelSubscriptions, error)
	GetSubscriptionsForChannel(ctx context.Context, channelKey keys.ChannelKey) ([]*chat.ChannelSubscription, error)
	GetSubscription(ctx context.Context, request *GetSubscription) (*chat.ChannelSubscription, error)
	GetMessage(ctx context.Context, messageKey keys.MessageKey) (*chat.Message, error)
	GetMessages(ctx context.Context, request *GetMessages) (*GetMessagesResponse, error)
	SaveMessage(ctx context.Context, request *chat.Message) error
	GetChannel(ctx context.Context, channelKey keys.ChannelKey) (*chat.Channel, error)
	CreateChannel(ctx context.Context, channel *chat.Channel) error
	CreateSubscription(ctx context.Context, key keys.ChannelSubscriptionKey, name string) (*chat.ChannelSubscription, error)
	DeleteSubscription(ctx context.Context, key keys.ChannelSubscriptionKey) error
}

type GetMessage struct {
	MessageKey keys.MessageKey
}

type GetMessageResponse struct {
	Message *chat.Message
}

type GetChannel struct {
	ChannelKey keys.ChannelKey
}

type GetChannelResponse struct {
	Channel *chat.Channel
}

type SaveMessageRequest struct {
	ChannelKey    keys.ChannelKey   `json:"channelKey"`
	Text          string            `json:"text"`
	Attachments   []chat.Attachment `json:"attachments"`
	Blocks        []chat.Block      `json:"blocks"`
	FromUser      keys.UserKey      `json:"fromUser"`
	FromUserName  string            `json:"toUser"`
	VisibleToUser *keys.UserKey     `json:"visibleToUser"`
}

type SaveMessageResponse struct {
	Message *chat.Message
}

type SendMessageToThreadResponse struct {
}

type GetMessages struct {
	Take    int
	Before  time.Time
	Channel keys.ChannelKey
	UserKey keys.UserKey
}

type GetMessagesResponse struct {
	Messages chat.Messages
	HasMore  bool
}

type GetSubscription struct {
	SubscriptionKey keys.ChannelSubscriptionKey
}

type GetSubscriptionResponse struct {
	Subscription *chat.ChannelSubscription
}

type GetSubscriptions struct {
	Take    int
	Skip    int
	UserKey keys.UserKey
}

func NewGetSubscriptions(userKey keys.UserKey, take int, skip int) *GetSubscriptions {
	return &GetSubscriptions{
		Take:    take,
		Skip:    skip,
		UserKey: userKey,
	}
}

type GetSubscriptionsResponse struct {
	Subscriptions chat.ChannelSubscriptions
}
