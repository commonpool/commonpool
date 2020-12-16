package model

import (
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
	"github.com/satori/go.uuid"
)

type OfferDecisionKey struct {
	OfferID uuid.UUID
	UserID  string
}

func NewOfferDecisionKey(offerKey OfferKey, userKey usermodel.UserKey) OfferDecisionKey {
	return OfferDecisionKey{
		OfferID: offerKey.ID,
		UserID:  userKey.String(),
	}
}

func (o *OfferDecisionKey) GetUserKey() usermodel.UserKey {
	return usermodel.NewUserKey(o.UserID)
}

func (o *OfferDecisionKey) GetOfferKey() OfferKey {
	return NewOfferKey(o.OfferID)
}
