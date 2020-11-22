package web

import (
	"github.com/commonpool/backend/model"
	"time"
)

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

type SendOfferRequest struct {
	Offer SendOfferPayload `json:"offer" validate:"required"`
}

type SendOfferPayload struct {
	Items []SendOfferPayloadItem `json:"items" validate:"min=1"`
}

type SendOfferPayloadItem struct {
	From          string              `json:"from" validate:"required,uuid"`
	To            string              `json:"to" validate:"required,uuid"`
	Type          model.OfferItemType `json:"type" validate:"required,min=0,max=1"`
	ResourceId    *string             `json:"resourceId" validate:"required,uuid"`
	TimeInSeconds *int64              `json:"timeInSeconds"`
}
