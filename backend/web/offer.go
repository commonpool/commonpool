package web

import (
	"github.com/commonpool/backend/trading"
	"time"
)

type Offer struct {
	ID             string              `json:"id"`
	CreatedAt      time.Time           `json:"createdAt"`
	CompletedAt    *time.Time          `json:"completedAt"`
	Status         trading.OfferStatus `json:"status"`
	AuthorID       string              `json:"authorId"`
	AuthorUsername string              `json:"authorUsername"`
	Items          []OfferItem         `json:"items"`
	Decisions      []OfferDecision     `json:"decisions"`
	Message        string              `json:"message"`
}

type OfferItem struct {
	ID            string                `json:"id"`
	FromUserID    string                `json:"fromUserId"`
	ToUserID      string                `json:"toUserId"`
	Type          trading.OfferItemType `json:"type"`
	ResourceId    *string               `json:"resourceId"`
	TimeInSeconds *int64                `json:"timeInSeconds"`
}

type OfferDecision struct {
	OfferID  string           `json:"offerId"`
	UserID   string           `json:"userId"`
	Decision trading.Decision `json:"decision"`
}

type GetOfferResponse struct {
	Offer Offer `json:"offer"`
}
type GetOffersResponse struct {
	Offers []Offer `json:"offers"`
}

type SendOfferRequest struct {
	Offer SendOfferPayload `json:"offer" validate:"required"`
}

type SendOfferPayload struct {
	Items   []SendOfferPayloadItem `json:"items" validate:"min=1"`
	Message string                 `json:"message"`
}

type SendOfferPayloadItem struct {
	From          string                `json:"from" validate:"required,uuid"`
	To            string                `json:"to" validate:"required,uuid"`
	Type          trading.OfferItemType `json:"type" validate:"required,min=0,max=1"`
	ResourceId    *string               `json:"resourceId" validate:"required,uuid"`
	TimeInSeconds *int64                `json:"timeInSeconds"`
}

func NewSendOfferPayloadItemForResource(from string, to string, resourceId string) *SendOfferPayloadItem {
	return &SendOfferPayloadItem{
		From:          from,
		To:            to,
		Type:          trading.ResourceItem,
		ResourceId:    &resourceId,
		TimeInSeconds: nil,
	}
}

func NewSendOfferPayloadItemForTime(from string, to string, time int64) *SendOfferPayloadItem {
	return &SendOfferPayloadItem{
		From:          from,
		To:            to,
		Type:          trading.TimeItem,
		ResourceId:    nil,
		TimeInSeconds: &time,
	}
}
