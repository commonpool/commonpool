package chat

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type Store interface {
	GetSubscriptionsForUser(ctx context.Context, request *GetSubscriptions) (*ChannelSubscriptions, error)
	GetSubscriptionsForChannel(ctx context.Context, channelKey keys.ChannelKey) ([]*ChannelSubscription, error)
	GetSubscription(ctx context.Context, request *GetSubscription) (*ChannelSubscription, error)
	GetMessage(ctx context.Context, messageKey keys.MessageKey) (*Message, error)
	GetMessages(ctx context.Context, request *GetMessages) (*GetMessagesResponse, error)
	SaveMessage(ctx context.Context, request *Message) error
	GetChannel(ctx context.Context, channelKey keys.ChannelKey) (*Channel, error)
	CreateChannel(ctx context.Context, channel *Channel) error
	CreateSubscription(ctx context.Context, key keys.ChannelSubscriptionKey, name string) (*ChannelSubscription, error)
	DeleteSubscription(ctx context.Context, key keys.ChannelSubscriptionKey) error
}

type GetMessage struct {
	MessageKey keys.MessageKey
}

type GetMessageResponse struct {
	Message *Message
}

type GetChannel struct {
	ChannelKey keys.ChannelKey
}

type GetChannelResponse struct {
	Channel *Channel
}

type SaveMessageRequest struct {
	ChannelKey    keys.ChannelKey `json:"channelKey"`
	Text          string          `json:"text"`
	Attachments   []Attachment    `json:"attachments"`
	Blocks        []Block         `json:"blocks"`
	FromUser      keys.UserKey    `json:"fromUser"`
	FromUserName  string          `json:"toUser"`
	VisibleToUser *keys.UserKey   `json:"visibleToUser"`
}

type SaveMessageResponse struct {
	Message *Message
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
	Messages Messages
	HasMore  bool
}

type GetSubscription struct {
	SubscriptionKey keys.ChannelSubscriptionKey
}

type GetSubscriptionResponse struct {
	Subscription *ChannelSubscription
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
	Subscriptions ChannelSubscriptions
}
