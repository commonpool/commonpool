package model

import (
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
	"time"
)

type Offer struct {
	Key            OfferKey
	GroupKey       groupmodel.GroupKey
	CreatedByKey   usermodel.UserKey
	Status         OfferStatus
	CreatedAt      time.Time
	ExpirationTime *time.Time
	CompletedAt    *time.Time
	Message        string
}

func NewOffer(offerKey OfferKey, groupKey groupmodel.GroupKey, author usermodel.UserKey, message string, expiration *time.Time) *Offer {
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

func (o *Offer) GetAuthorKey() usermodel.UserKey {
	return o.CreatedByKey
}

func (o *Offer) IsPending() bool {
	return o.Status == PendingOffer
}
