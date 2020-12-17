package chat

import (
	ctx "context"
	"github.com/commonpool/backend/pkg/keys"
	"golang.org/x/net/context"
	"time"
)

type Service interface {
	GetMessages(ctx context.Context, channel keys.ChannelKey, before time.Time, take int) (*GetMessagesResponse, error)
	GetSubscriptionsForUser(ctx context.Context, take int, skip int) (*ChannelSubscriptions, error)
	GetChannel(ctx context.Context, channelKey keys.ChannelKey) (*Channel, error)
	GetMessage(ctx context.Context, messageKey keys.MessageKey) (*Message, error)
	CreateChannel(ctx ctx.Context, channelKey keys.ChannelKey, channelType ChannelType) (*Channel, error)
	SubscribeToChannel(ctx ctx.Context, channelSubscriptionKey keys.ChannelSubscriptionKey, name string) (*ChannelSubscription, error)
	UnsubscribeFromChannel(ctx context.Context, channelSubscriptionKey keys.ChannelSubscriptionKey) error
	DeleteGroupChannel(ctx ctx.Context, request *DeleteGroupChannel) (*DeleteGroupChannelResponse, error)
	SendConversationMessage(ctx ctx.Context, request *SendConversationMessage) (*SendConversationMessageResponse, error)
	SendMessage(ctx context.Context, message *Message) error
	SendGroupMessage(ctx ctx.Context, request *SendGroupMessage) (*SendGroupMessageResponse, error)
	CreateUserExchange(ctx context.Context, userKey keys.UserKey) (string, error)
	GetUserExchangeName(ctx context.Context, userKey keys.UserKey) string
}
type GetOrCreateConversationChannelResponse struct {
	Channel *Channel
}

type DeleteGroupChannel struct {
	GroupKey keys.GroupKey
}

type DeleteGroupChannelResponse struct {
	Channel *Channel
}

type SendConversationMessage struct {
	FromUserKey          keys.UserKey
	FromUserName         string
	ToUserKeys           *keys.UserKeys
	Text                 string
	Blocks               []Block
	Attachments          []Attachment
	OnlyVisibleToUserKey *keys.UserKey
}

type SendConversationMessageResponse struct {
	Message *Message
}

func NewSendConversationMessage(
	fromUserKey keys.UserKey,
	fromUserName string,
	toUserKeys *keys.UserKeys,
	text string,
	blocks []Block,
	attachments []Attachment,
	onlyVisibleToUserKey *keys.UserKey) *SendConversationMessage {
	return &SendConversationMessage{
		FromUserKey:          fromUserKey,
		FromUserName:         fromUserName,
		ToUserKeys:           toUserKeys,
		Text:                 text,
		Blocks:               blocks,
		Attachments:          attachments,
		OnlyVisibleToUserKey: onlyVisibleToUserKey,
	}
}

type SendGroupMessage struct {
	GroupKey             keys.GroupKey
	FromUserKey          keys.UserKey
	FromUserName         string
	Text                 string
	Blocks               []Block
	Attachments          []Attachment
	OnlyVisibleToUserKey *keys.UserKey
}

type SendGroupMessageResponse struct {
	Channel *Channel
}

func NewSendGroupMessage(groupKey keys.GroupKey, fromUserKey keys.UserKey, fromUserName string, text string, blocks []Block, attachments []Attachment, onlyVisibleToUserKey *keys.UserKey) *SendGroupMessage {
	return &SendGroupMessage{
		GroupKey:             groupKey,
		FromUserKey:          fromUserKey,
		FromUserName:         fromUserName,
		Text:                 text,
		Blocks:               blocks,
		Attachments:          attachments,
		OnlyVisibleToUserKey: onlyVisibleToUserKey,
	}
}
