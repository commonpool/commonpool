package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	"gorm.io/gorm"
)

type GetResource struct {
	db *gorm.DB
}

func NewGetResource(db *gorm.DB) *GetResource {
	return &GetResource{db: db}
}

func (q *GetResource) Get(ctx context.Context, resourceKey keys.ResourceKey) (*readmodel.ResourceReadModel, error) {
	var resource readmodel.ResourceReadModel
	qry := q.db.Model(&readmodel.ResourceReadModel{}).Where("resource_key = ?", resourceKey.String()).Find(&resource)
	if qry.Error != nil {
		return nil, qry.Error
	}
	if qry.RowsAffected == 0 {
		return nil, exceptions.ErrResourceNotFound
	}
	return &resource, nil
}
