package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/group/readmodels"
	"github.com/commonpool/backend/pkg/keys"
	"gorm.io/gorm"
)

type GetGroup struct {
	db *gorm.DB
}

func NewGetGroupReadModel(db *gorm.DB) *GetGroup {
	return &GetGroup{
		db: db,
	}
}

func (q *GetGroup) Get(ctx context.Context, groupKey keys.GroupKey) (*readmodels.GroupReadModel, error) {
	var rm readmodels.GroupReadModel
	if err := q.db.Model(&readmodels.GroupReadModel{}).Where("group_key = ?", groupKey.String()).Find(&rm).Error; err != nil {
		return nil, err
	}
	return &rm, nil
}
