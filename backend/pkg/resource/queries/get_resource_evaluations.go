package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	"gorm.io/gorm"
)

type GetResourceEvaluations struct {
	db *gorm.DB
}

func NewGetResourceEvaluations(db *gorm.DB) *GetResourceEvaluations {
	return &GetResourceEvaluations{db: db}
}

func (q *GetResourceEvaluations) Get(ctx context.Context, resourceKey keys.ResourceKey) ([]*readmodel.ResourceEvaluationReadModel, error) {
	var result []*readmodel.ResourceEvaluationReadModel
	if err := q.db.Model(&readmodel.ResourceEvaluationReadModel{}).
		Where("resource_key = ?", resourceKey).
		Order("evaluated_at asc").
		Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
