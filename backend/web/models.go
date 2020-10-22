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
	Id                string `json:"id"`
	WithUsername      string `json:"withUsername"`
	WithUsernameId    string `json:"withUsernameId"`
	LastChars         string `json:"lastChars"`
	HasUnreadMessages bool   `json:"hasUnreadMessages"`
}

type Message struct {
	ID         string    `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	SenderId   string    `json:"senderId"`
	ReceiverId string    `json:"receiverId"`
	ThreadId   string    `json:"threadId"`
	Content    string    `json:"content"`
}

type GetThreadMessagesResponse struct {
	Messages []Message `json:"messages"`
}
