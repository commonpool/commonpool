package model

import (
	"github.com/commonpool/backend/pkg/chat/chatmodel"
	"time"
)

type Subscription struct {
	ChannelID           string                `json:"channelId"`
	UserID              string                `json:"userId"`
	HasUnreadMessages   bool                  `json:"hasUnreadMessages"`
	CreatedAt           time.Time             `json:"createdAt"`
	UpdatedAt           time.Time             `json:"updatedAt"`
	LastMessageAt       time.Time             `json:"lastMessageAt"`
	LastTimeRead        time.Time             `json:"lastTimeRead"`
	LastMessageChars    string                `json:"lastMessageChars"`
	LastMessageUserId   string                `json:"lastMessageUserId"`
	LastMessageUserName string                `json:"lastMessageUsername"`
	Name                string                `json:"name"`
	Type                chatmodel.ChannelType `json:"type"`
}

func MapSubscription(channel *chatmodel.Channel, subscription *chatmodel.ChannelSubscription) *Subscription {
	return &Subscription{
		ChannelID:           channel.Key.String(),
		UserID:              subscription.UserKey.String(),
		HasUnreadMessages:   subscription.LastMessageAt.After(subscription.LastTimeRead),
		CreatedAt:           subscription.CreatedAt,
		UpdatedAt:           subscription.UpdatedAt,
		LastMessageAt:       subscription.LastMessageAt,
		LastTimeRead:        subscription.LastTimeRead,
		LastMessageChars:    subscription.LastMessageChars,
		LastMessageUserId:   subscription.LastMessageUserKey.String(),
		LastMessageUserName: subscription.LastMessageUserName,
		Name:                subscription.Name,
		Type:                channel.Type,
	}
}

type Message struct {
	ID             string                   `json:"id"`
	ChannelID      string                   `json:"channelId"`
	MessageType    chatmodel.MessageType    `json:"messageType"`
	MessageSubType chatmodel.MessageSubType `json:"messageSubType"`
	SentById       string                   `json:"sentById"`
	SentByUsername string                   `json:"sentByUsername"`
	SentAt         time.Time                `json:"sentAt"`
	Text           string                   `json:"text"`
	Blocks         []chatmodel.Block        `json:"blocks"`
	Attachments    []chatmodel.Attachment   `json:"attachments"`
	VisibleToUser  *string                  `json:"visibleToUser"`
}

func MapMessage(message *chatmodel.Message) *Message {
	var visibleToUser *string = nil
	if message.VisibleToUser != nil {
		visibleToUserStr := message.VisibleToUser.String()
		visibleToUser = &visibleToUserStr
	}
	return &Message{
		ID:             message.Key.String(),
		ChannelID:      message.ChannelKey.String(),
		MessageType:    message.MessageType,
		MessageSubType: message.MessageSubType,
		SentById:       message.SentBy.UserKey.String(),
		SentByUsername: message.SentBy.Username,
		SentAt:         message.SentAt,
		Text:           message.Text,
		Blocks:         message.Blocks,
		Attachments:    message.Attachments,
		VisibleToUser:  visibleToUser,
	}
}

type GetLatestSubscriptionsResponse struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

type InquireAboutResourceRequest struct {
	Message string `json:"message"`
}

type SendMessageRequest struct {
	Message string `json:"message,omitempty" validate:"notblank,required,min=1,max=2000"`
}

type GetLatestMessageThreadsResponse struct {
	Messages []Message `json:"messages"`
}

type GetTopicMessagesResponse struct {
	Messages []*Message `json:"messages"`
}

type InteractionMessage struct {
	Payload InteractionCallbackPayload `json:"payload"`
}

type InteractionPayloadType string

const (
	BlockActions InteractionPayloadType = "block_actions"
)

type ElementState struct {
	Type            chatmodel.ElementType    `json:"type,omitempty"`
	SelectedDate    *string                  `json:"selectedDate,omitempty"`
	SelectedTime    *string                  `json:"selectedTime,omitempty"`
	Value           *string                  `json:"value,omitempty"`
	SelectedOption  *chatmodel.OptionObject  `json:"selectedOption,omitempty"`
	SelectedOptions []chatmodel.OptionObject `json:"selectedOptions,omitempty"`
}

type SubmitAction struct {
	ElementState
	BlockID  string `json:"blockId,omitempty"`
	ActionID string `json:"actionId,omitempty"`
}

type Action struct {
	SubmitAction
	ActionTimestamp time.Time `json:"actionTimestamp,omitempty"`
}

type SubmitInteractionPayload struct {
	MessageID string                             `json:"messageId,omitempty"`
	State     map[string]map[string]ElementState `json:"state,omitempty"`
	Actions   []SubmitAction                     `json:"actions,omitempty"`
}

type SubmitInteractionRequest struct {
	Payload SubmitInteractionPayload `json:"payload"`
}

type InteractionCallback struct {
	Payload InteractionCallbackPayload `json:"payload"`
	Token   string                     `json:"token"`
}

type InteractionPayloadUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type InteractionCallbackPayload struct {
	Type        InteractionPayloadType             `json:"type"`
	User        InteractionPayloadUser             `json:"user"`
	TriggerId   string                             `json:"triggerId"`
	ResponseURL string                             `json:"responseUrl"`
	Message     *Message                           `json:"message"`
	Actions     []Action                           `json:"actions"`
	State       map[string]map[string]ElementState `json:"state"`
}
