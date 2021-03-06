package queries

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/listeners"
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

func (q *GetOfferKeyForOfferItemKey) Get(offerItemKey keys.OfferItemKey) (keys.OfferKey, error) {
	var readModel listeners.OfferItemReadModel
	if err := q.db.
		Model(&listeners.OfferItemReadModel{}).
		Where("id = ?", offerItemKey.String()).
		First(&readModel).
		Error; err != nil {
		return keys.OfferKey{}, err
	}
	return keys.ParseOfferKey(readModel.OfferID)
}
