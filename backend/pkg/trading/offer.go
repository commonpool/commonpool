package trading

import (
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type Offer struct {
	Key            keys.OfferKey `json:"id"`
	GroupKey       keys.GroupKey `json:"groupId"`
	CreatedByKey   keys.UserKey  `json:"createdById"`
	Status         OfferStatus   `json:"status"`
	CreatedAt      time.Time     `json:"createdAt"`
	ExpirationTime time.Time     `json:"expirationTime"`
	CompletedAt    time.Time     `json:"completedAt"`
	Message        string        `json:"message"`
}

func NewOffer(offerKey keys.OfferKey, groupKey keys.GroupKey, author keys.UserKey, message string, expiration time.Time) *Offer {
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

func (o *Offer) GetKey() keys.OfferKey {
	return o.Key
}

func (o *Offer) GetAuthorKey() keys.UserKey {
	return o.CreatedByKey
}

func (o *Offer) IsPending() bool {
	return o.Status == PendingOffer
}
