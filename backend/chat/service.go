package chat

import (
	ctx "context"
	"github.com/commonpool/backend/model"
	"golang.org/x/net/context"
	"time"
)

type Service interface {
	NotifyUserInterestedAboutResource(ctx ctx.Context, request *NotifyUserInterestedAboutResource) (*NotifyUserInterestedAboutResourceResponse, error)
	CreateChannel(ctx ctx.Context, channelKey model.ChannelKey, channelType ChannelType) (*Channel, error)
	SubscribeToChannel(ctx ctx.Context, channelSubscriptionKey model.ChannelSubscriptionKey, name string) (*ChannelSubscription, error)
	UnsubscribeFromChannel(ctx context.Context, channelSubscriptionKey model.ChannelSubscriptionKey) error
	DeleteGroupChannel(ctx ctx.Context, request *DeleteGroupChannel) (*DeleteGroupChannelResponse, error)
	SendConversationMessage(ctx ctx.Context, request *SendConversationMessage) (*SendConversationMessageResponse, error)
	SendChannelMessage(ctx context.Context, channelKey model.ChannelKey, message string) (*Message, error)
	SendGroupMessage(ctx ctx.Context, request *SendGroupMessage) (*SendGroupMessageResponse, error)
	CreateUserExchange(ctx context.Context, userKey model.UserKey) (string, error)
	GetUserExchangeName(ctx context.Context, userKey model.UserKey) string
}

type NotifyUserInterestedAboutResource struct {
	InterestedUser model.UserKey
	ResourceKey    model.ResourceKey
	Message        string
}
type NotifyUserInterestedAboutResourceResponse struct {
	ChannelKey model.ChannelKey
}

func NewNotifyUserInterestedAboutResource(interestedUser model.UserKey, resourceKey model.ResourceKey, message string) *NotifyUserInterestedAboutResource {
	return &NotifyUserInterestedAboutResource{
		InterestedUser: interestedUser,
		ResourceKey:    resourceKey,
		Message:        message,
	}
}

type NotifyUserOffer struct {
	OfferingUser model.UserKey
	OfferKey     model.OfferKey
}
type NotifyUserOfferResponse struct {
	Message *Message
}

func NewNotifyUserOffer(offeringUser model.UserKey, offerKey model.OfferKey) *NotifyUserOffer {
	return &NotifyUserOffer{
		OfferingUser: offeringUser,
		OfferKey:     offerKey,
	}
}

type NotifyOfferAccepted struct {
	AcceptingUser model.UserKey
	Offer         model.OfferKey
}

type NotifyOfferAcceptedResponse struct {
	Message *Message
}

func NewNotifyOfferAccepted(acceptingUser model.UserKey, offerKey model.OfferKey) *NotifyOfferAccepted {
	return &NotifyOfferAccepted{
		AcceptingUser: acceptingUser,
		Offer:         offerKey,
	}
}

type NotifyOfferDeclined struct {
	DecliningUser model.UserKey
	Offer         model.OfferKey
}
type NotifyOfferDeclinedResponse struct {
	Message *Message
}

func NewNotifyOfferDeclined(decliningUser model.UserKey, offerKey model.OfferKey) *NotifyOfferDeclined {
	return &NotifyOfferDeclined{
		DecliningUser: decliningUser,
		Offer:         offerKey,
	}
}

type NotifyGroupJoined struct {
	UserKey  model.UserKey
	GroupKey model.GroupKey
}

type NotifyGroupJoinedResponse struct {
	Message *Message
}

func NewNotifyGroupJoined(userKey model.UserKey, groupKey model.GroupKey) *NotifyGroupJoined {
	return &NotifyGroupJoined{
		UserKey:  userKey,
		GroupKey: groupKey,
	}
}

type NotifyGroupLeft struct {
	UserKey  model.UserKey
	GroupKey model.GroupKey
}

type NotifyGroupLeftResponse struct {
	Message *Message
}

func NewNotifyGroupLeft(userKey model.UserKey, groupKey model.GroupKey) *NotifyGroupLeft {
	return &NotifyGroupLeft{
		UserKey:  userKey,
		GroupKey: groupKey,
	}
}

type NotifyCreditsTransferred struct {
	FromUserKey model.UserKey
	ToUserKey   model.GroupKey
	Amount      time.Time
}

type NotifyCreditsTransferredResponse struct {
	Message *Message
}

func NewNotifyCreditsTransferred(fromUserKey model.UserKey, toUserKey model.GroupKey, amount time.Time) *NotifyCreditsTransferred {
	return &NotifyCreditsTransferred{
		FromUserKey: fromUserKey,
		ToUserKey:   toUserKey,
		Amount:      amount,
	}
}

type GetOrCreateConversationChannel struct {
	UserKeys *model.UserKeys
}

type GetOrCreateConversationChannelResponse struct {
	Channel *Channel
}

func NewGetOrCreateConversationChannel(userKeys *model.UserKeys) *GetOrCreateConversationChannel {
	return &GetOrCreateConversationChannel{
		UserKeys: userKeys,
	}
}

type DeleteGroupChannel struct {
	GroupKey model.GroupKey
}

type DeleteGroupChannelResponse struct {
	Channel *Channel
}

func NewDeleteGroupChannel(groupKey model.GroupKey) *DeleteGroupChannel {
	return &DeleteGroupChannel{
		GroupKey: groupKey,
	}
}

type SendConversationMessage struct {
	FromUserKey          model.UserKey
	FromUserName         string
	ToUserKeys           *model.UserKeys
	Text                 string
	Blocks               []Block
	Attachments          []Attachment
	OnlyVisibleToUserKey *model.UserKey
}

type SendConversationMessageResponse struct {
	Message *Message
}

func NewSendConversationMessage(
	fromUserKey model.UserKey,
	fromUserName string,
	toUserKeys *model.UserKeys,
	text string,
	blocks []Block,
	attachments []Attachment,
	onlyVisibleToUserKey *model.UserKey) *SendConversationMessage {
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
	GroupKey             model.GroupKey
	FromUserKey          model.UserKey
	FromUserName         string
	Text                 string
	Blocks               []Block
	Attachments          []Attachment
	OnlyVisibleToUserKey *model.UserKey
}

type SendGroupMessageResponse struct {
	Channel *Channel
}

func NewSendGroupMessage(groupKey model.GroupKey, fromUserKey model.UserKey, fromUserName string, text string, blocks []Block, attachments []Attachment, onlyVisibleToUserKey *model.UserKey) *SendGroupMessage {
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
