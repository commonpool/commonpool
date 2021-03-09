package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	"gorm.io/gorm"
	"strings"
)

type GetResourcesByKeys struct {
	db *gorm.DB
}

func NewGetResourcesByKeys(db *gorm.DB) *GetResourcesByKeys {
	return &GetResourcesByKeys{db: db}
}

func (q *GetResourcesByKeys) Get(ctx context.Context, resourceKeys *keys.ResourceKeys) ([]*readmodel.ResourceReadModel, error) {
	if resourceKeys.IsEmpty() {
		return []*readmodel.ResourceReadModel{}, nil
	}
	var sb strings.Builder
	var params []interface{}
	sb.WriteString("resource_key in (")
	keyCount := len(resourceKeys.Items)
	for i, item := range resourceKeys.Items {
		sb.WriteString("?")
		if i < keyCount-1 {
			sb.WriteString(",")
		}
		params = append(params, item)
	}
	sb.WriteString(")")
	var result []*readmodel.DbResourceReadModel
	qry := q.db.Model(&readmodel.DbResourceReadModel{}).Where(sb.String(), params...).Find(&result)
	if qry.Error != nil {
		return nil, qry.Error
	}
	if int(qry.RowsAffected) != keyCount {
		return nil, exceptions.ErrNotFoundf("Some resource keys could not be found")
	}
	return mapResourceReadModels(result), nil
}
