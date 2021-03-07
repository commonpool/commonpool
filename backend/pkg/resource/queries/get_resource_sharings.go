package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	"gorm.io/gorm"
)

type GetResourceSharings struct {
	db *gorm.DB
}

func NewGetResourceSharings(db *gorm.DB) *GetResourceSharings {
	return &GetResourceSharings{db: db}
}

func (q *GetResourceSharings) Get(ctx context.Context, resourceKey keys.ResourceKey) ([]*readmodel.ResourceSharingReadModel, error) {
	var sharings []*readmodel.ResourceSharingReadModel
	if err := q.db.Find(&sharings, "resource_key = ?", resourceKey.String()).Error; err != nil {
		return nil, err
	}
	return sharings, nil
}
