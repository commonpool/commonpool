package service

import (
	ctx "context"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/chat/store"
	"github.com/commonpool/backend/pkg/keys"
	"golang.org/x/net/context"
	"time"
)

type Service interface {
	GetMessages(ctx context.Context, channel keys.ChannelKey, before time.Time, take int) (*store.GetMessagesResponse, error)
	GetSubscriptionsForUser(ctx context.Context, take int, skip int) (*chat.ChannelSubscriptions, error)
	GetChannel(ctx context.Context, channelKey keys.ChannelKey) (*chat.Channel, error)
	GetMessage(ctx context.Context, messageKey keys.MessageKey) (*chat.Message, error)
	CreateChannel(ctx ctx.Context, channelKey keys.ChannelKey, channelType chat.ChannelType) (*chat.Channel, error)
	SubscribeToChannel(ctx ctx.Context, channelSubscriptionKey keys.ChannelSubscriptionKey, name string) (*chat.ChannelSubscription, error)
	UnsubscribeFromChannel(ctx context.Context, channelSubscriptionKey keys.ChannelSubscriptionKey) error
	DeleteGroupChannel(ctx ctx.Context, request *DeleteGroupChannel) (*DeleteGroupChannelResponse, error)
	SendConversationMessage(ctx ctx.Context, request *SendConversationMessage) (*SendConversationMessageResponse, error)
	SendMessage(ctx context.Context, message *chat.Message) error
	SendGroupMessage(ctx ctx.Context, request *SendGroupMessage) (*SendGroupMessageResponse, error)
	CreateUserExchange(ctx context.Context, userKey keys.UserKey) (string, error)
	GetUserExchangeName(ctx context.Context, userKey keys.UserKey) string
}
type GetOrCreateConversationChannelResponse struct {
	Channel *chat.Channel
}

type DeleteGroupChannel struct {
	GroupKey keys.GroupKey
}

type DeleteGroupChannelResponse struct {
	Channel *chat.Channel
}

type SendConversationMessage struct {
	FromUserKey          keys.UserKey
	FromUserName         string
	ToUserKeys           *keys.UserKeys
	Text                 string
	Blocks               []chat.Block
	Attachments          []chat.Attachment
	OnlyVisibleToUserKey *keys.UserKey
}

type SendConversationMessageResponse struct {
	Message *chat.Message
}

func NewSendConversationMessage(
	fromUserKey keys.UserKey,
	fromUserName string,
	toUserKeys *keys.UserKeys,
	text string,
	blocks []chat.Block,
	attachments []chat.Attachment,
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
	Blocks               []chat.Block
	Attachments          []chat.Attachment
	OnlyVisibleToUserKey *keys.UserKey
}

type SendGroupMessageResponse struct {
	Channel *chat.Channel
}

func NewSendGroupMessage(groupKey keys.GroupKey, fromUserKey keys.UserKey, fromUserName string, text string, blocks []chat.Block, attachments []chat.Attachment, onlyVisibleToUserKey *keys.UserKey) *SendGroupMessage {
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
