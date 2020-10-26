package model

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type OfferStatus int

const (
	PendingOffer OfferStatus = iota
	AcceptedOffer
	CanceledOffer
	DeclinedOffer
	ExpiredOffer
)

type Offer struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key"`
	AuthorID       string
	Status         OfferStatus
	CreatedAt      time.Time
	ExpirationTime *time.Time
	// Completion time (when either accepted, expired or declined)
	CompletedAt *time.Time
}

func NewOffer(offerKey OfferKey, author UserKey, expiration *time.Time) Offer {
	return Offer{
		ID:             offerKey.ID,
		AuthorID:       author.String(),
		Status:         PendingOffer,
		ExpirationTime: expiration,
	}
}

func (o *Offer) GetKey() OfferKey {
	return NewOfferKey(o.ID)
}

func (o *Offer) GetAuthorKey() UserKey {
	return NewUserKey(o.AuthorID)
}
