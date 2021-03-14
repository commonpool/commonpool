package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"gorm.io/gorm"
)

type GetResourcesEvaluations struct {
	db *gorm.DB
}

func NewGetResourcesEvaluations(db *gorm.DB) *GetResourceEvaluations {
	return &GetResourceEvaluations{db: db}
}

func (q *GetResourcesEvaluations) Get(ctx context.Context, resourceKey keys.ResourceKey) {

}
