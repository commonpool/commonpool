package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/domain"
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
	var resource readmodel.DbResourceReadModel
	qry := q.db.Model(&readmodel.DbResourceReadModel{}).Where("resource_key = ?", resourceKey.String()).Find(&resource)
	if qry.Error != nil {
		return nil, qry.Error
	}
	if qry.RowsAffected == 0 {
		return nil, exceptions.ErrResourceNotFound
	}
	return mapResourceReadModel(&resource), nil

}

func mapResourceReadModel(resource *readmodel.DbResourceReadModel) *readmodel.ResourceReadModel {
	return &readmodel.ResourceReadModel{
		ResourceReadModelBase: resource.ResourceReadModelBase,
		ResourceInfo: domain.ResourceInfo{
			ResourceInfoBase: resource.ResourceInfoBase,
		},
	}
}
func mapResourceReadModels(resources []*readmodel.DbResourceReadModel) []*readmodel.ResourceReadModel {
	var result []*readmodel.ResourceReadModel
	for _, resource := range resources {
		result = append(result, &readmodel.ResourceReadModel{
			ResourceReadModelBase: resource.ResourceReadModelBase,
			ResourceInfo: domain.ResourceInfo{
				ResourceInfoBase: resource.ResourceInfoBase,
			},
		})
	}
	return result
}
