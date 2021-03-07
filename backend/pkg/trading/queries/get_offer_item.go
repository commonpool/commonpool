package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	keys2 "github.com/commonpool/backend/pkg/trading/readmodels"
	"gorm.io/gorm"
)

type GetOfferItem struct {
	db *gorm.DB
}

func NewGetOfferItem(db *gorm.DB) *GetOfferItem {
	return &GetOfferItem{db: db}
}

func (q *GetOfferItem) Get(ctx context.Context, offerItemKey keys.OfferItemKey) (*keys2.OfferItemReadModel2, error) {
	var result keys2.OfferItemReadModel
	err := q.db.Model(&keys2.OfferItemReadModel{}).Find(&result, "id = ?", offerItemKey.String()).Error

	cache := NewReadModelCache()
	cache.processOfferItem(&result)
	if err := cache.retrieve(q.db); err != nil {
		return nil, err
	}

	return mapOfferItem(&result, cache), err
}
