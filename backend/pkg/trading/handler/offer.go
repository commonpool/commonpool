package handler

import (
	tradingdomain "github.com/commonpool/backend/pkg/trading/domain"
	"time"
)

type Offer struct {
	ID             string                    `json:"id"`
	CreatedAt      time.Time                 `json:"createdAt"`
	CompletedAt    time.Time                 `json:"completedAt"`
	Status         tradingdomain.OfferStatus `json:"status"`
	AuthorID       string                    `json:"authorId"`
	AuthorUsername string                    `json:"authorUsername"`
	Items          []*OfferItem              `json:"items"`
	Message        string                    `json:"message"`
}

type OfferResponse struct {
	Offer *Offer `json:"offer"`
}
