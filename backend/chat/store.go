package chat

import (
	"context"
	"github.com/commonpool/backend/model"
	"time"
)

type Store interface {
	GetSubscriptionsForUser(ctx context.Context, request *GetSubscriptions) (*ChannelSubscriptions, error)
	GetSubscriptionsForChannel(ctx context.Context, channelKey model.ChannelKey) ([]ChannelSubscription, error)
	GetSubscription(ctx context.Context, request *GetSubscription) (*ChannelSubscription, error)
	GetMessage(ctx context.Context, messageKey model.MessageKey) (*Message, error)
	GetMessages(ctx context.Context, request *GetMessages) (*GetMessagesResponse, error)
	SaveMessage(ctx context.Context, request *SaveMessageRequest) (*Message, error)
	GetChannel(ctx context.Context, channelKey model.ChannelKey) (*Channel, error)
	CreateChannel(ctx context.Context, channel *Channel) error
	CreateSubscription(ctx context.Context, key model.ChannelSubscriptionKey, name string) (*ChannelSubscription, error)
	DeleteSubscription(ctx context.Context, key model.ChannelSubscriptionKey) error
}

type GetMessage struct {
	MessageKey model.MessageKey
}

type GetMessageResponse struct {
	Message *Message
}

type GetChannel struct {
	ChannelKey model.ChannelKey
}

type GetChannelResponse struct {
	Channel *Channel
}

type SaveMessageRequest struct {
	ChannelKey    model.ChannelKey `json:"channelKey"`
	Text          string           `json:"text"`
	Attachments   []Attachment     `json:"attachments"`
	Blocks        []Block          `json:"blocks"`
	FromUser      model.UserKey    `json:"fromUser"`
	FromUserName  string           `json:"toUser"`
	VisibleToUser *model.UserKey   `json:"visibleToUser"`
}

type SaveMessageResponse struct {
	Message *Message
}

type SendMessageToThreadResponse struct {
}

func NewSaveMessageRequest(
	topicKey model.ChannelKey,
	fromUser model.UserKey,
	fromUserName string,
	text string,
	blocks []Block,
	attachments []Attachment,
	visibleToUserOnly *model.UserKey,
) *SaveMessageRequest {
	return &SaveMessageRequest{
		ChannelKey:    topicKey,
		Text:          text,
		Attachments:   attachments,
		Blocks:        blocks,
		FromUser:      fromUser,
		FromUserName:  fromUserName,
		VisibleToUser: visibleToUserOnly,
	}
}

type GetMessages struct {
	Take    int
	Before  time.Time
	Channel model.ChannelKey
	UserKey model.UserKey
}

type GetMessagesResponse struct {
	Messages Messages
	HasMore  bool
}

type GetSubscription struct {
	SubscriptionKey model.ChannelSubscriptionKey
}

type GetSubscriptionResponse struct {
	Subscription *ChannelSubscription
}

type GetSubscriptions struct {
	Take    int
	Skip    int
	UserKey model.UserKey
}

func NewGetSubscriptions(userKey model.UserKey, take int, skip int) *GetSubscriptions {
	return &GetSubscriptions{
		Take:    take,
		Skip:    skip,
		UserKey: userKey,
	}
}

type GetSubscriptionsResponse struct {
	Subscriptions ChannelSubscriptions
}
