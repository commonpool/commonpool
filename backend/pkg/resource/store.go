package resource

import (
	"context"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
)

type Store interface {
	GetByKey(ctx context.Context, getResourceByKeyQuery *GetResourceByKeyQuery) (*GetResourceByKeyResponse, error)
	GetByKeys(ctx context.Context, resourceKeys *resourcemodel.ResourceKeys) (*GetResourceByKeysResponse, error)
	Search(ctx context.Context, searchResourcesQuery *SearchResourcesQuery) (*SearchResourcesResponse, error)
	Delete(ctx context.Context, resourceKey resourcemodel.ResourceKey) error
	Create(ctx context.Context, createResourceQuery *CreateResourceQuery) error
	Update(ctx context.Context, updateResourceQuery *UpdateResourceQuery) error
}

type GetResourceByKeysQuery struct {
	ResourceKeys []resourcemodel.ResourceKey
}

type SearchResourcesQuery struct {
	Type            *resourcemodel.Type
	SubType         *resourcemodel.SubType
	Query           *string
	Skip            int
	Take            int
	CreatedBy       string
	SharedWithGroup *groupmodel.GroupKey
}

func NewSearchResourcesQuery(query *string, resourceType *resourcemodel.Type, resourceSubType *resourcemodel.SubType, skip int, take int, createdBy string, sharedWithGroup *groupmodel.GroupKey) *SearchResourcesQuery {
	return &SearchResourcesQuery{
		Type:            resourceType,
		SubType:         resourceSubType,
		Query:           query,
		Skip:            skip,
		Take:            take,
		CreatedBy:       createdBy,
		SharedWithGroup: sharedWithGroup,
	}
}

type SearchResourcesResponse struct {
	Resources  *resourcemodel.Resources
	Sharings   *resourcemodel.Sharings
	TotalCount int
	Skip       int
	Take       int
}

type GetResourceByKeyQuery struct {
	ResourceKey resourcemodel.ResourceKey
}

func NewGetResourceByKeyQuery(resourceKey resourcemodel.ResourceKey) *GetResourceByKeyQuery {
	return &GetResourceByKeyQuery{
		ResourceKey: resourceKey,
	}
}

type GetResourceByKeyResponse struct {
	Resource *resourcemodel.Resource
	Sharings *resourcemodel.Sharings
	Claims   *resourcemodel.Claims
}

type GetResourceByKeysResponse struct {
	Resources *resourcemodel.Resources
	Sharings  *resourcemodel.Sharings
	Claims    *resourcemodel.Claims
}

type CreateResourceQuery struct {
	Resource   *resourcemodel.Resource
	SharedWith *groupmodel.GroupKeys
}

func NewCreateResourceQuery(resource *resourcemodel.Resource, sharedWith *groupmodel.GroupKeys) *CreateResourceQuery {
	return &CreateResourceQuery{
		Resource:   resource,
		SharedWith: sharedWith,
	}
}

type UpdateResourceQuery struct {
	Resource   *resourcemodel.Resource
	SharedWith *groupmodel.GroupKeys
}

func NewUpdateResourceQuery(resource *resourcemodel.Resource, sharedWith *groupmodel.GroupKeys) *UpdateResourceQuery {
	return &UpdateResourceQuery{
		Resource:   resource,
		SharedWith: sharedWith,
	}
}
