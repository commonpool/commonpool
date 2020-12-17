package keys

import (
	"github.com/satori/go.uuid"
)

type OfferDecisionKey struct {
	OfferID uuid.UUID
	UserID  string
}

func NewOfferDecisionKey(offerKey OfferKey, userKey UserKey) OfferDecisionKey {
	return OfferDecisionKey{
		OfferID: offerKey.ID,
		UserID:  userKey.String(),
	}
}

func (o *OfferDecisionKey) GetUserKey() UserKey {
	return NewUserKey(o.UserID)
}

func (o *OfferDecisionKey) GetOfferKey() OfferKey {
	return NewOfferKey(o.OfferID)
}
