package web

import (
	"github.com/commonpool/backend/model"
	"time"
)

type Thread struct {
	TopicID             string    `json:"topicId"`
	RecipientID         string    `json:"recipientId"`
	LastChars           string    `json:"lastChars"`
	HasUnreadMessages   bool      `json:"hasUnreadMessages"`
	LastMessageAt       time.Time `json:"lastMessageAt"`
	LastMessageUsername string    `json:"lastMessageUsername"`
	LastMessageUserId   string    `json:"lastMessageUserId"`
	Title               string    `json:"title"`
}

type Message struct {
	ID             string               `json:"id"`
	TopicID        string               `json:"topicId"`
	MessageType    model.MessageType    `json:"messageType"`
	MessageSubType model.MessageSubType `json:"messageSubType"`
	UserID         string               `json:"userId"`
	BotID          string               `json:"botId"`
	SentAt         time.Time            `json:"sentAt"`
	Text           string               `json:"text"`
	Blocks         []model.Block        `json:"blocks"`
	Attachments    []model.Attachment   `json:"attachments"`
	IsPersonal     bool                 `json:"isPersonal"`
	SentBy         string               `json:"sentBy"`
	SentByUsername string               `json:"sentByUsername"`
}

type GetLatestThreadsResponse struct {
	Threads []Thread `json:"threads"`
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
	Payload InteractionPayload `json:"payload"`
}

type InteractionPayloadType string

const (
	BlockActions InteractionPayloadType = "block_actions"
)

type SubmitInteractionMessage struct {
	ID             string               `json:"id"`
	TopicID        string               `json:"topicId"`
	Blocks         []model.Block        `json:"blocks"`
	Attachments    []model.Attachment   `json:"attachments"`
}

type SubmitInteractionRequest struct {
	Payload SubmitInteractionPayload `json:"payload"`
}

type SubmitInteractionPayload struct {

}

type InteractionPayloadUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type InteractionPayloadAction struct {
	BlockID  string `json:"blockId"`
	ActionID string `json:"actionId"`
	Value    string `json:"value"`
}

type InteractionPayload struct {
	Type        InteractionPayloadType       `json:"type"`
	TriggertId  string                       `json:"triggerId"`
	ResponseURL string                       `json:"responseUrl"`
	User        InteractionPayloadUser       `json:"user"`
	Message     Message                      `json:"message"`
	Actions     []InteractionPayloadAction   `json:"actions"`
	State       map[string]map[string]string `json:"state"`
}
