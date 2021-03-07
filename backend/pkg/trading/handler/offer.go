package handler

import (
	"time"
)

type Offer struct {
	ID             string       `json:"id"`
	CreatedAt      time.Time    `json:"createdAt"`
	CompletedAt    time.Time    `json:"completedAt"`
	Status         string       `json:"status"`
	AuthorID       string       `json:"authorId"`
	AuthorUsername string       `json:"authorUsername"`
	Items          []*OfferItem `json:"items"`
	Message        string       `json:"message"`
}

type OfferResponse struct {
	Offer *Offer `json:"offer"`
}
