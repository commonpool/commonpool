package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	keys2 "github.com/commonpool/backend/pkg/trading/readmodels"
	"gorm.io/gorm"
)

type GetOfferKeyForOfferItemKey struct {
	db *gorm.DB
}

func NewGetOfferKeyForOfferItemKey(db *gorm.DB) *GetOfferKeyForOfferItemKey {
	return &GetOfferKeyForOfferItemKey{
		db: db,
	}
}

func (q *GetOfferKeyForOfferItemKey) Get(ctx context.Context, offerItemKey keys.OfferItemKey) (keys.OfferKey, error) {
	var readModel keys2.OfferItemReadModel
	if err := q.db.
		Model(&keys2.OfferItemReadModel{}).
		Where("id = ?", offerItemKey.String()).
		First(&readModel).
		Error; err != nil {
		return keys.OfferKey{}, err
	}
	return readModel.OfferKey, nil
}
