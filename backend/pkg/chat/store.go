package chat

import (
	"context"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
	"time"
)

type Store interface {
	GetSubscriptionsForUser(ctx context.Context, request *GetSubscriptions) (*ChannelSubscriptions, error)
	GetSubscriptionsForChannel(ctx context.Context, channelKey ChannelKey) ([]*ChannelSubscription, error)
	GetSubscription(ctx context.Context, request *GetSubscription) (*ChannelSubscription, error)
	GetMessage(ctx context.Context, messageKey MessageKey) (*Message, error)
	GetMessages(ctx context.Context, request *GetMessages) (*GetMessagesResponse, error)
	SaveMessage(ctx context.Context, request *Message) error
	GetChannel(ctx context.Context, channelKey ChannelKey) (*Channel, error)
	CreateChannel(ctx context.Context, channel *Channel) error
	CreateSubscription(ctx context.Context, key ChannelSubscriptionKey, name string) (*ChannelSubscription, error)
	DeleteSubscription(ctx context.Context, key ChannelSubscriptionKey) error
}

type GetMessage struct {
	MessageKey MessageKey
}

type GetMessageResponse struct {
	Message *Message
}

type GetChannel struct {
	ChannelKey ChannelKey
}

type GetChannelResponse struct {
	Channel *Channel
}

type SaveMessageRequest struct {
	ChannelKey    ChannelKey         `json:"channelKey"`
	Text          string             `json:"text"`
	Attachments   []Attachment       `json:"attachments"`
	Blocks        []Block            `json:"blocks"`
	FromUser      usermodel.UserKey  `json:"fromUser"`
	FromUserName  string             `json:"toUser"`
	VisibleToUser *usermodel.UserKey `json:"visibleToUser"`
}

type SaveMessageResponse struct {
	Message *Message
}

type SendMessageToThreadResponse struct {
}

type GetMessages struct {
	Take    int
	Before  time.Time
	Channel ChannelKey
	UserKey usermodel.UserKey
}

type GetMessagesResponse struct {
	Messages Messages
	HasMore  bool
}

type GetSubscription struct {
	SubscriptionKey ChannelSubscriptionKey
}

type GetSubscriptionResponse struct {
	Subscription *ChannelSubscription
}

type GetSubscriptions struct {
	Take    int
	Skip    int
	UserKey usermodel.UserKey
}

func NewGetSubscriptions(userKey usermodel.UserKey, take int, skip int) *GetSubscriptions {
	return &GetSubscriptions{
		Take:    take,
		Skip:    skip,
		UserKey: userKey,
	}
}

type GetSubscriptionsResponse struct {
	Subscriptions ChannelSubscriptions
}
