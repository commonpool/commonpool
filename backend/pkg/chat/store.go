package chat

import (
	"context"
	chatmodel "github.com/commonpool/backend/pkg/chat/model"
	usermodel "github.com/commonpool/backend/pkg/user/model"
	"time"
)

type Store interface {
	GetSubscriptionsForUser(ctx context.Context, request *GetSubscriptions) (*chatmodel.ChannelSubscriptions, error)
	GetSubscriptionsForChannel(ctx context.Context, channelKey chatmodel.ChannelKey) ([]*chatmodel.ChannelSubscription, error)
	GetSubscription(ctx context.Context, request *GetSubscription) (*chatmodel.ChannelSubscription, error)
	GetMessage(ctx context.Context, messageKey chatmodel.MessageKey) (*chatmodel.Message, error)
	GetMessages(ctx context.Context, request *GetMessages) (*GetMessagesResponse, error)
	SaveMessage(ctx context.Context, request *chatmodel.Message) error
	GetChannel(ctx context.Context, channelKey chatmodel.ChannelKey) (*chatmodel.Channel, error)
	CreateChannel(ctx context.Context, channel *chatmodel.Channel) error
	CreateSubscription(ctx context.Context, key chatmodel.ChannelSubscriptionKey, name string) (*chatmodel.ChannelSubscription, error)
	DeleteSubscription(ctx context.Context, key chatmodel.ChannelSubscriptionKey) error
}

type GetMessage struct {
	MessageKey chatmodel.MessageKey
}

type GetMessageResponse struct {
	Message *chatmodel.Message
}

type GetChannel struct {
	ChannelKey chatmodel.ChannelKey
}

type GetChannelResponse struct {
	Channel *chatmodel.Channel
}

type SaveMessageRequest struct {
	ChannelKey    chatmodel.ChannelKey   `json:"channelKey"`
	Text          string                 `json:"text"`
	Attachments   []chatmodel.Attachment `json:"attachments"`
	Blocks        []chatmodel.Block      `json:"blocks"`
	FromUser      usermodel.UserKey      `json:"fromUser"`
	FromUserName  string                 `json:"toUser"`
	VisibleToUser *usermodel.UserKey     `json:"visibleToUser"`
}

type SaveMessageResponse struct {
	Message *chatmodel.Message
}

type SendMessageToThreadResponse struct {
}

type GetMessages struct {
	Take    int
	Before  time.Time
	Channel chatmodel.ChannelKey
	UserKey usermodel.UserKey
}

type GetMessagesResponse struct {
	Messages chatmodel.Messages
	HasMore  bool
}

type GetSubscription struct {
	SubscriptionKey chatmodel.ChannelSubscriptionKey
}

type GetSubscriptionResponse struct {
	Subscription *chatmodel.ChannelSubscription
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
	Subscriptions chatmodel.ChannelSubscriptions
}
