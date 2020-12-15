package chat

import (
	ctx "context"
	chatmodel "github.com/commonpool/backend/pkg/chat/model"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
	usermodel "github.com/commonpool/backend/pkg/user/model"
	"golang.org/x/net/context"
	"time"
)

type Service interface {
	GetMessages(ctx context.Context, channel chatmodel.ChannelKey, before time.Time, take int) (*GetMessagesResponse, error)
	GetSubscriptionsForUser(ctx context.Context, take int, skip int) (*chatmodel.ChannelSubscriptions, error)
	GetChannel(ctx context.Context, channelKey chatmodel.ChannelKey) (*chatmodel.Channel, error)
	GetMessage(ctx context.Context, messageKey chatmodel.MessageKey) (*chatmodel.Message, error)
	NotifyUserInterestedAboutResource(ctx ctx.Context, request *NotifyUserInterestedAboutResource) (*NotifyUserInterestedAboutResourceResponse, error)
	CreateChannel(ctx ctx.Context, channelKey chatmodel.ChannelKey, channelType chatmodel.ChannelType) (*chatmodel.Channel, error)
	SubscribeToChannel(ctx ctx.Context, channelSubscriptionKey chatmodel.ChannelSubscriptionKey, name string) (*chatmodel.ChannelSubscription, error)
	UnsubscribeFromChannel(ctx context.Context, channelSubscriptionKey chatmodel.ChannelSubscriptionKey) error
	DeleteGroupChannel(ctx ctx.Context, request *DeleteGroupChannel) (*DeleteGroupChannelResponse, error)
	SendConversationMessage(ctx ctx.Context, request *SendConversationMessage) (*SendConversationMessageResponse, error)
	SendMessage(ctx context.Context, message *chatmodel.Message) error
	SendGroupMessage(ctx ctx.Context, request *SendGroupMessage) (*SendGroupMessageResponse, error)
	CreateUserExchange(ctx context.Context, userKey usermodel.UserKey) (string, error)
	GetUserExchangeName(ctx context.Context, userKey usermodel.UserKey) string
	GetUserLink(userKey usermodel.UserKey) string
	GetGroupLink(groupKey groupmodel.GroupKey) string
	GetResourceLink(resource resourcemodel.ResourceKey) string
}

type NotifyUserInterestedAboutResource struct {
	InterestedUser usermodel.UserKey
	ResourceKey    resourcemodel.ResourceKey
	Message        string
}
type NotifyUserInterestedAboutResourceResponse struct {
	ChannelKey chatmodel.ChannelKey
}

func NewNotifyUserInterestedAboutResource(interestedUser usermodel.UserKey, resourceKey resourcemodel.ResourceKey, message string) *NotifyUserInterestedAboutResource {
	return &NotifyUserInterestedAboutResource{
		InterestedUser: interestedUser,
		ResourceKey:    resourceKey,
		Message:        message,
	}
}

type GetOrCreateConversationChannelResponse struct {
	Channel *chatmodel.Channel
}

type DeleteGroupChannel struct {
	GroupKey groupmodel.GroupKey
}

type DeleteGroupChannelResponse struct {
	Channel *chatmodel.Channel
}

type SendConversationMessage struct {
	FromUserKey          usermodel.UserKey
	FromUserName         string
	ToUserKeys           *usermodel.UserKeys
	Text                 string
	Blocks               []chatmodel.Block
	Attachments          []chatmodel.Attachment
	OnlyVisibleToUserKey *usermodel.UserKey
}

type SendConversationMessageResponse struct {
	Message *chatmodel.Message
}

func NewSendConversationMessage(
	fromUserKey usermodel.UserKey,
	fromUserName string,
	toUserKeys *usermodel.UserKeys,
	text string,
	blocks []chatmodel.Block,
	attachments []chatmodel.Attachment,
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
	GroupKey             groupmodel.GroupKey
	FromUserKey          usermodel.UserKey
	FromUserName         string
	Text                 string
	Blocks               []chatmodel.Block
	Attachments          []chatmodel.Attachment
	OnlyVisibleToUserKey *usermodel.UserKey
}

type SendGroupMessageResponse struct {
	Channel *chatmodel.Channel
}

func NewSendGroupMessage(groupKey groupmodel.GroupKey, fromUserKey usermodel.UserKey, fromUserName string, text string, blocks []chatmodel.Block, attachments []chatmodel.Attachment, onlyVisibleToUserKey *usermodel.UserKey) *SendGroupMessage {
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
