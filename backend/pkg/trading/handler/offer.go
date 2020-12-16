package handler

import (
	"github.com/commonpool/backend/pkg/trading"
	"time"
)

type Offer struct {
	ID             string              `json:"id"`
	CreatedAt      time.Time           `json:"createdAt"`
	CompletedAt    *time.Time          `json:"completedAt"`
	Status         trading.OfferStatus `json:"status"`
	AuthorID       string              `json:"authorId"`
	AuthorUsername string              `json:"authorUsername"`
	Items          []*OfferItem        `json:"items"`
	Message        string              `json:"message"`
}

type OfferResponse struct {
	Offer *Offer `json:"offer"`
}
