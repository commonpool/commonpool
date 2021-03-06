package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/listeners"
	"gorm.io/gorm"
)

type GetOfferResult struct {
	Offer      *listeners.OfferReadModel
	OfferItems []*listeners.OfferItemReadModel
}

type GetOffer struct {
	db *gorm.DB
}

func NewGetOffer(db *gorm.DB) *GetOffer {
	return &GetOffer{
		db: db,
	}
}

func (q *GetOffer) Get(ctx context.Context, offerKey keys.OfferKey) (*GetOfferResult, error) {

	var offer listeners.OfferReadModel
	if err := q.db.Model(&listeners.OfferReadModel{}).Find(&offer, "id = ?", offerKey.ID.String()).Error; err != nil {
		return nil, err
	}

	var offerItems []*listeners.OfferItemReadModel
	if err := q.db.Model(&listeners.OfferItemReadModel{}).Find(&offerItems, "offer_id = ?", offerKey.ID.String()).Error; err != nil {
		return nil, err
	}

	return &GetOfferResult{
		Offer:      &offer,
		OfferItems: offerItems,
	}, nil

}
