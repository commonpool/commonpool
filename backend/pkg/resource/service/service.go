package service

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource"
)

type ResourceService struct {
	resourceStore resource.Store
}

var _ resource.Service = &ResourceService{}

func NewResourceService(resourceStore resource.Store) *ResourceService {
	return &ResourceService{
		resourceStore: resourceStore,
	}
}

func (r ResourceService) GetResourcesByKeys(ctx context.Context, resourceKeys *keys.ResourceKeys) (*resource.Resources, error) {
	response, err := r.resourceStore.GetByKeys(ctx, resourceKeys)
	if err != nil {
		return nil, err
	}
	return response.Resources, nil
}

func (r ResourceService) GetByKey(ctx context.Context, query *resource.GetResourceByKeyQuery) (*resource.GetResourceByKeyResponse, error) {
	return r.resourceStore.GetByKey(ctx, query)
}

func (r ResourceService) Search(ctx context.Context, query *resource.SearchResourcesQuery) (*resource.SearchResourcesResponse, error) {
	return r.resourceStore.Search(ctx, query)
}

func (r ResourceService) Create(ctx context.Context, query *resource.CreateResourceQuery) error {
	return r.resourceStore.Create(ctx, query)
}

func (r ResourceService) Update(ctx context.Context, query *resource.UpdateResourceQuery) error {
	return r.resourceStore.Update(ctx, query)
}
