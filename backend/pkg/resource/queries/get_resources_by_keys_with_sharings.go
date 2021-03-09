package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"strings"
)

type GetResourcesByKeysWithSharings struct {
	db *gorm.DB
}

func NewGetResourcesByKeysWithSharings(db *gorm.DB) *GetResourcesByKeysWithSharings {
	return &GetResourcesByKeysWithSharings{db: db}
}

func (q *GetResourcesByKeysWithSharings) Get(ctx context.Context, resourceKeys *keys.ResourceKeys) ([]*readmodel.ResourceWithSharingsReadModel, error) {

	if resourceKeys.IsEmpty() {
		return []*readmodel.ResourceWithSharingsReadModel{}, nil
	}

	g, ctx := errgroup.WithContext(ctx)

	var sb strings.Builder
	var params []interface{}
	sb.WriteString("resource_key in (")
	for i, resourceKey := range resourceKeys.Items {
		sb.WriteString("?")
		params = append(params, resourceKey)
		if i < resourceKeys.Count()-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString(")")

	var dbResources []*readmodel.DbResourceReadModel
	var sharings []*readmodel.ResourceSharingReadModel
	g.Go(func() error {
		return q.db.Model(&readmodel.DbResourceReadModel{}).Where(sb.String(), params...).Find(&dbResources).Error
	})

	g.Go(func() error {
		return q.db.Model(&readmodel.ResourceSharingReadModel{}).Where(sb.String(), params...).Find(&sharings).Error
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return mapResourcesWithSharings(mapResourceReadModels(dbResources), sharings), nil
}
