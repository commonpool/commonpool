package resource

import (
	"context"
)

type Service interface {
	GetResourcesByKeys(ctx context.Context, resourceKeys *ResourceKeys) (*Resources, error)
	GetByKey(ctx context.Context, query *GetResourceByKeyQuery) (*GetResourceByKeyResponse, error)
	Search(ctx context.Context, query *SearchResourcesQuery) (*SearchResourcesResponse, error)
	Create(ctx context.Context, query *CreateResourceQuery) error
	Update(ctx context.Context, query *UpdateResourceQuery) error
}
