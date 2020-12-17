package trading

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/satori/go.uuid"
)

type OfferDecisionKey struct {
	OfferID uuid.UUID
	UserID  string
}

func NewOfferDecisionKey(offerKey OfferKey, userKey keys.UserKey) OfferDecisionKey {
	return OfferDecisionKey{
		OfferID: offerKey.ID,
		UserID:  userKey.String(),
	}
}

func (o *OfferDecisionKey) GetUserKey() keys.UserKey {
	return keys.NewUserKey(o.UserID)
}

func (o *OfferDecisionKey) GetOfferKey() OfferKey {
	return NewOfferKey(o.OfferID)
}
