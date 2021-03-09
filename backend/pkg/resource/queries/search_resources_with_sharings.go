package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	"gorm.io/gorm"
)

type SearchResourcesWithSharings struct {
	db                   *gorm.DB
	searchResources      *SearchResources
	getResourcesSharings *GetResourcesSharings
}

func NewSearchResourcesWithSharings(db *gorm.DB, searchResources *SearchResources, getResourcesSharings *GetResourcesSharings) *SearchResourcesWithSharings {
	return &SearchResourcesWithSharings{
		db:                   db,
		searchResources:      searchResources,
		getResourcesSharings: getResourcesSharings,
	}
}

func (q *SearchResourcesWithSharings) Get(ctx context.Context, query *SearchResourcesQuery) ([]*readmodel.ResourceWithSharingsReadModel, error) {

	resources, err := q.searchResources.Get(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(resources) == 0 {
		return []*readmodel.ResourceWithSharingsReadModel{}, nil
	}

	var resourceKeys []keys.ResourceKey
	for _, resource := range resources {
		resourceKeys = append(resourceKeys, resource.ResourceKey)
	}

	sharings, err := q.getResourcesSharings.Get(ctx, keys.NewResourceKeys(resourceKeys))
	if err != nil {
		return nil, err
	}

	groupedSharings := map[keys.ResourceKey][]*readmodel.ResourceSharingReadModel{}
	for _, sharing := range sharings {
		groupedSharings[sharing.ResourceKey] = append(groupedSharings[sharing.ResourceKey], sharing)
	}

	var result []*readmodel.ResourceWithSharingsReadModel
	for _, resource := range resources {
		sharings, ok := groupedSharings[resource.ResourceKey]
		if !ok {
			sharings = []*readmodel.ResourceSharingReadModel{}
		}
		result = append(result, &readmodel.ResourceWithSharingsReadModel{
			ResourceReadModel: *resource,
			Sharings:          sharings,
		})
	}

	return result, nil

}
