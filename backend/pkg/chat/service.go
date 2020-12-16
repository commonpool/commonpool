package chat

import (
	ctx "context"
	"github.com/commonpool/backend/pkg/group"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
	"golang.org/x/net/context"
	"time"
)

type Service interface {
	GetMessages(ctx context.Context, channel ChannelKey, before time.Time, take int) (*GetMessagesResponse, error)
	GetSubscriptionsForUser(ctx context.Context, take int, skip int) (*ChannelSubscriptions, error)
	GetChannel(ctx context.Context, channelKey ChannelKey) (*Channel, error)
	GetMessage(ctx context.Context, messageKey MessageKey) (*Message, error)
	CreateChannel(ctx ctx.Context, channelKey ChannelKey, channelType ChannelType) (*Channel, error)
	SubscribeToChannel(ctx ctx.Context, channelSubscriptionKey ChannelSubscriptionKey, name string) (*ChannelSubscription, error)
	UnsubscribeFromChannel(ctx context.Context, channelSubscriptionKey ChannelSubscriptionKey) error
	DeleteGroupChannel(ctx ctx.Context, request *DeleteGroupChannel) (*DeleteGroupChannelResponse, error)
	SendConversationMessage(ctx ctx.Context, request *SendConversationMessage) (*SendConversationMessageResponse, error)
	SendMessage(ctx context.Context, message *Message) error
	SendGroupMessage(ctx ctx.Context, request *SendGroupMessage) (*SendGroupMessageResponse, error)
	CreateUserExchange(ctx context.Context, userKey usermodel.UserKey) (string, error)
	GetUserExchangeName(ctx context.Context, userKey usermodel.UserKey) string
}
type GetOrCreateConversationChannelResponse struct {
	Channel *Channel
}

type DeleteGroupChannel struct {
	GroupKey group.GroupKey
}

type DeleteGroupChannelResponse struct {
	Channel *Channel
}

type SendConversationMessage struct {
	FromUserKey          usermodel.UserKey
	FromUserName         string
	ToUserKeys           *usermodel.UserKeys
	Text                 string
	Blocks               []Block
	Attachments          []Attachment
	OnlyVisibleToUserKey *usermodel.UserKey
}

type SendConversationMessageResponse struct {
	Message *Message
}

func NewSendConversationMessage(
	fromUserKey usermodel.UserKey,
	fromUserName string,
	toUserKeys *usermodel.UserKeys,
	text string,
	blocks []Block,
	attachments []Attachment,
	onlyVisibleToUserKey *usermodel.UserKey) *SendConversationMessage {
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
	GroupKey             group.GroupKey
	FromUserKey          usermodel.UserKey
	FromUserName         string
	Text                 string
	Blocks               []Block
	Attachments          []Attachment
	OnlyVisibleToUserKey *usermodel.UserKey
}

type SendGroupMessageResponse struct {
	Channel *Channel
}

func NewSendGroupMessage(groupKey group.GroupKey, fromUserKey usermodel.UserKey, fromUserName string, text string, blocks []Block, attachments []Attachment, onlyVisibleToUserKey *usermodel.UserKey) *SendGroupMessage {
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
