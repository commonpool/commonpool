package resource

import (
	"context"
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
)

type Service interface {
	GetResourcesByKeys(ctx context.Context, resourceKeys *resourcemodel.ResourceKeys) (*resourcemodel.Resources, error)
	GetByKey(ctx context.Context, query *GetResourceByKeyQuery) (*GetResourceByKeyResponse, error)
	Search(ctx context.Context, query *SearchResourcesQuery) (*SearchResourcesResponse, error)
	Create(ctx context.Context, query *CreateResourceQuery) error
	Update(ctx context.Context, query *UpdateResourceQuery) error
}
