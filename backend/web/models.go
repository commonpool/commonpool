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

type UsersInfoResponse struct {
	Users []UserInfoResponse `json:"users"`
	Take  int                `json:"take"`
	Skip  int                `json:"skip"`
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
	Title               string    `json:"title"`
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

type SendOfferRequest struct {
	Offer SendOfferPayload `json:"offer"`
}

type SendOfferPayload struct {
	Items []SendOfferPayloadItem `json:"items"`
}

type SendOfferPayloadItem struct {
	From          string              `json:"from"`
	To            string              `json:"to"`
	Type          model.OfferItemType `json:"type"`
	ResourceId    *string             `json:"resourceId"`
	TimeInSeconds *int64              `json:"timeInSeconds"`
}

type Offer struct {
	ID             string            `json:"id"`
	CreatedAt      time.Time         `json:"createdAt"`
	CompletedAt    *time.Time        `json:"completedAt"`
	Status         model.OfferStatus `json:"status"`
	AuthorID       string            `json:"authorId"`
	AuthorUsername string            `json:"authorUsername"`
	Items          []OfferItem       `json:"items"`
	Decisions      []OfferDecision   `json:"decisions"`
}

type OfferItem struct {
	ID            string              `json:"id"`
	FromUserID    string              `json:"fromUserId"`
	ToUserID      string              `json:"toUserId"`
	Type          model.OfferItemType `json:"type"`
	ResourceId    string              `json:"resourceId"`
	TimeInSeconds int64               `json:"timeInSeconds"`
}

type OfferDecision struct {
	OfferID  string         `json:"offerId"`
	UserID   string         `json:"userId"`
	Decision model.Decision `json:"decision"`
}

type GetOfferResponse struct {
	Offer Offer `json:"offer"`
}
type GetOffersResponse struct {
	Offers []Offer `json:"offers"`
}
