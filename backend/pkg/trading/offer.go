package trading

import (
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type Offer struct {
	Key            OfferKey
	GroupKey       keys.GroupKey
	CreatedByKey   keys.UserKey
	Status         OfferStatus
	CreatedAt      time.Time
	ExpirationTime *time.Time
	CompletedAt    *time.Time
	Message        string
}

func NewOffer(offerKey OfferKey, groupKey keys.GroupKey, author keys.UserKey, message string, expiration *time.Time) *Offer {
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

func (o *Offer) GetKey() OfferKey {
	return o.Key
}

func (o *Offer) GetAuthorKey() keys.UserKey {
	return o.CreatedByKey
}

func (o *Offer) IsPending() bool {
	return o.Status == PendingOffer
}
