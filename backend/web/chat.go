package web

import (
	"github.com/commonpool/backend/chat"
	"time"
)

type Subscription struct {
	ChannelID           string           `json:"channelId"`
	UserID              string           `json:"userId"`
	HasUnreadMessages   bool             `json:"hasUnreadMessages"`
	CreatedAt           time.Time        `json:"createdAt"`
	UpdatedAt           time.Time        `json:"updatedAt"`
	LastMessageAt       time.Time        `json:"lastMessageAt"`
	LastTimeRead        time.Time        `json:"lastTimeRead"`
	LastMessageChars    string           `json:"lastMessageChars"`
	LastMessageUserId   string           `json:"lastMessageUserId"`
	LastMessageUserName string           `json:"lastMessageUsername"`
	Name                string           `json:"name"`
	Type                chat.ChannelType `json:"type"`
}

type Message struct {
	ID             string              `json:"id"`
	ChannelID      string              `json:"channelId"`
	MessageType    chat.MessageType    `json:"messageType"`
	MessageSubType chat.MessageSubType `json:"messageSubType"`
	SentById       string              `json:"sentById"`
	SentByUsername string              `json:"sentByUsername"`
	SentAt         time.Time           `json:"sentAt"`
	Text           string              `json:"text"`
	Blocks         []chat.Block        `json:"blocks"`
	Attachments    []chat.Attachment   `json:"attachments"`
	VisibleToUser  *string             `json:"visibleToUser"`
}

type GetLatestSubscriptionsResponse struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

type InquireAboutResourceRequest struct {
	Message string `json:"message"`
}

type SendMessageRequest struct {
	Message string `json:"message"`
}

type GetLatestMessageThreadsResponse struct {
	Messages []Message `json:"messages"`
}

type GetTopicMessagesResponse struct {
	Messages []Message `json:"messages"`
}

type InteractionMessage struct {
	Payload InteractionCallbackPayload `json:"payload"`
}

type InteractionPayloadType string

const (
	BlockActions InteractionPayloadType = "block_actions"
)

type ElementState struct {
	Type            chat.ElementType    `json:"type,omitempty"`
	SelectedDate    *string             `json:"selectedDate,omitempty"`
	SelectedTime    *string             `json:"selectedTime,omitempty"`
	Value           *string             `json:"value,omitempty"`
	SelectedOption  *chat.OptionObject  `json:"selectedOption,omitempty"`
	SelectedOptions []chat.OptionObject `json:"selectedOptions,omitempty"`
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
	Message     Message                            `json:"message"`
	Actions     []Action                           `json:"actions"`
	State       map[string]map[string]ElementState `json:"state"`
}
