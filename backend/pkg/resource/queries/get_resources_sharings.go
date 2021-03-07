package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	"gorm.io/gorm"
	"strings"
)

type GetResourcesSharings struct {
	db *gorm.DB
}

func NewGetResourcesSharings(db *gorm.DB) *GetResourcesSharings {
	return &GetResourcesSharings{db: db}
}

func (q *GetResourcesSharings) Get(ctx context.Context, resourceKeys *keys.ResourceKeys) ([]*readmodel.ResourceSharingReadModel, error) {
	if resourceKeys.IsEmpty() {
		return []*readmodel.ResourceSharingReadModel{}, nil
	}
	var sharings []*readmodel.ResourceSharingReadModel
	var sql = "resource_key in ("
	var params []interface{}
	var paramsStrs []string
	var visitedMap map[string]bool
	for _, item := range resourceKeys.Items {
		resourceKeyStr := item.String()
		if visitedMap[resourceKeyStr] {
			continue
		}
		params = append(params, resourceKeyStr)
		paramsStrs = append(paramsStrs, "?")
	}
	sql = "where resource_key in (" + strings.Join(paramsStrs, ",") + ")"
	if err := q.db.Find(&sharings, sql).Error; err != nil {
		return nil, err
	}
	return sharings, nil
}
