package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	groupreadmodels "github.com/commonpool/backend/pkg/trading/readmodels"
	"gorm.io/gorm"
)

type GetOffers struct {
	db                  *gorm.DB
	getOfferKeysForUser *GetOfferKeysForUser
}

func NewGetOffers(db *gorm.DB, getOfferKeysForUser *GetOfferKeysForUser) *GetOffers {
	return &GetOffers{db: db, getOfferKeysForUser: getOfferKeysForUser}
}

func (q *GetOffers) Get(ctx context.Context, userKey keys.UserKey) ([]*groupreadmodels.OfferReadModel, error) {

	offerKeys, err := q.getOfferKeysForUser.Get(ctx, userKey)
	if err != nil {
		return nil, err
	}

	return getOffers(ctx, offerKeys, q.db)

}
