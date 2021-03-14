package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	"gorm.io/gorm"
)

type GetUserResourceEvaluation struct {
	db *gorm.DB
}

func NewGetUserResourceEvaluation(db *gorm.DB) *GetUserResourceEvaluation {
	return &GetUserResourceEvaluation{db: db}
}

func (q *GetUserResourceEvaluation) Get(ctx context.Context, resourceKey keys.ResourceKey, userKey keys.UserKey) ([]*readmodel.ResourceEvaluationReadModel, error) {
	var result []*readmodel.ResourceEvaluationReadModel
	if err := q.db.Model(&readmodel.ResourceEvaluationReadModel{}).
		Where("resource_key = ? and evaluated_by = ?", resourceKey, userKey).
		Order("evaluated_at desc").
		Find(&result).Error; err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return []*readmodel.ResourceEvaluationReadModel{}, nil
	}

	var latestEvaluationId *string

	latestUserEvaluations := []*readmodel.ResourceEvaluationReadModel{}
	for _, model := range result {
		if latestEvaluationId == nil {
			latestEvaluationId = &model.EvaluationID
		} else {
			if model.EvaluationID != *latestEvaluationId {
				break
			}
		}
		latestUserEvaluations = append(latestUserEvaluations, model)
	}

	return latestUserEvaluations, nil
}
