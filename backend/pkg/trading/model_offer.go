package trading

import (
	"github.com/commonpool/backend/model"
	"time"
)

type Offer struct {
	Key            model.OfferKey
	GroupKey       model.GroupKey
	CreatedByKey   model.UserKey
	Status         OfferStatus
	CreatedAt      time.Time
	ExpirationTime *time.Time
	CompletedAt    *time.Time
	Message        string
}

func NewOffer(offerKey model.OfferKey, groupKey model.GroupKey, author model.UserKey, message string, expiration *time.Time) *Offer {
	return &Offer{
		Key:            offerKey,
		GroupKey:       groupKey,
		CreatedByKey:   author,
		Status:         PendingOffer,
		ExpirationTime: expiration,
		Message:        message,
		CreatedAt:      time.Now().UTC(),
	}
}

func (o *Offer) GetKey() model.OfferKey {
	return o.Key
}

func (o *Offer) GetAuthorKey() model.UserKey {
	return o.CreatedByKey
}

func (o *Offer) IsPending() bool {
	return o.Status == PendingOffer
}
