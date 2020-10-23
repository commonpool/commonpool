package web

import (
	"github.com/commonpool/backend/model"
	"time"
)

type SearchResourcesResponse struct {
	TotalCount int        `json:"totalCount"`
	Take       int        `json:"take"`
	Skip       int        `json:"skip"`
	Resources  []Resource `json:"resources"`
}

type GetResourceResponse struct {
	Resource Resource `json:"resource"`
}

type Resource struct {
	Id               string             `json:"id"`
	Summary          string             `json:"summary"`
	Description      string             `json:"description"`
	Type             model.ResourceType `json:"type"`
	CreatedAt        time.Time          `json:"createdAt"`
	CreatedBy        string             `json:"createdBy"`
	CreatedById      string             `json:"createdById"`
	ValueInHoursFrom int                `json:"valueInHoursFrom"`
	ValueInHoursTo   int                `json:"valueInHoursTo"`
}

type CreateResourceResponse struct {
	Resource Resource `json:"resource"`
}

type CreateResourceRequest struct {
	Resource CreateResourcePayload `json:"resource"`
}

type CreateResourcePayload struct {
	Summary          string             `json:"summary"`
	Description      string             `json:"description"`
	Type             model.ResourceType `json:"type"`
	ValueInHoursFrom int                `json:"valueInHoursFrom"`
	ValueInHoursTo   int                `json:"valueInHoursTo"`
}

type UpdateResourceRequest struct {
	Resource UpdateResourcePayload `json:"resource"`
}

type UpdateResourcePayload struct {
	Summary          string             `json:"summary"`
	Description      string             `json:"description"`
	Type             model.ResourceType `json:"type"`
	ValueInHoursFrom int                `json:"valueInHoursFrom"`
	ValueInHoursTo   int                `json:"valueInHoursTo"`
}

type UpdateResourceResponse struct {
	Resource Resource `json:"resource"`
}

type UserAuthResponse struct {
	IsAuthenticated bool   `json:"isAuthenticated"`
	Username        string `json:"username"`
	Id              string `json:"id"`
}

type UserInfoResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type GetLatestThreadsResponse struct {
	Threads []Thread `json:"threads"`
}

type Thread struct {
	TopicID             string    `json:"topicId"`
	RecipientID         string    `json:"recipientId"`
	LastChars           string    `json:"lastChars"`
	HasUnreadMessages   bool      `json:"hasUnreadMessages"`
	LastMessageAt       time.Time `json:"lastMessageAt"`
	LastMessageUsername string    `json:"lastMessageUsername"`
	LastMessageUserId   string    `json:"lastMessageUserId"`
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

type Message struct {
	ID             string    `json:"id"`
	TopicID        string    `json:"topicId"`
	SentBy         string    `json:"sentBy"`
	SentByUsername string    `json:"sentByUsername"`
	SentByMe       bool      `json:"sentByMe"`
	SentAt         time.Time `json:"sentAt"`
	Content        string    `json:"content"`
}

type GetTopicMessagesResponse struct {
	Messages []Message `json:"messages"`
}
