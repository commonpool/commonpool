package resource

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
)

type Store interface {
	GetByKey(ctx context.Context, getResourceByKeyQuery *GetResourceByKeyQuery) (*GetResourceByKeyResponse, error)
	GetByKeys(ctx context.Context, resourceKeys *ResourceKeys) (*GetResourceByKeysResponse, error)
	Search(ctx context.Context, searchResourcesQuery *SearchResourcesQuery) (*SearchResourcesResponse, error)
	Delete(ctx context.Context, resourceKey keys.ResourceKey) error
	Create(ctx context.Context, createResourceQuery *CreateResourceQuery) error
	Update(ctx context.Context, updateResourceQuery *UpdateResourceQuery) error
}

type GetResourceByKeysQuery struct {
	ResourceKeys []keys.ResourceKey
}

type SearchResourcesQuery struct {
	Type            *Type
	SubType         *SubType
	Query           *string
	Skip            int
	Take            int
	CreatedBy       string
	SharedWithGroup *keys.GroupKey
}

func NewSearchResourcesQuery(query *string, resourceType *Type, resourceSubType *SubType, skip int, take int, createdBy string, sharedWithGroup *keys.GroupKey) *SearchResourcesQuery {
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
	Resources  *Resources
	Sharings   *Sharings
	TotalCount int
	Skip       int
	Take       int
}

type GetResourceByKeyQuery struct {
	ResourceKey keys.ResourceKey
}

func NewGetResourceByKeyQuery(resourceKey keys.ResourceKey) *GetResourceByKeyQuery {
	return &GetResourceByKeyQuery{
		ResourceKey: resourceKey,
	}
}

type GetResourceByKeyResponse struct {
	Resource *Resource
	Sharings *Sharings
	Claims   *Claims
}

type GetResourceByKeysResponse struct {
	Resources *Resources
	Sharings  *Sharings
	Claims    *Claims
}

type CreateResourceQuery struct {
	Resource   *Resource
	SharedWith *keys.GroupKeys
}

func NewCreateResourceQuery(resource *Resource, sharedWith *keys.GroupKeys) *CreateResourceQuery {
	return &CreateResourceQuery{
		Resource:   resource,
		SharedWith: sharedWith,
	}
}

type UpdateResourceQuery struct {
	Resource   *Resource
	SharedWith *keys.GroupKeys
}

func NewUpdateResourceQuery(resource *Resource, sharedWith *keys.GroupKeys) *UpdateResourceQuery {
	return &UpdateResourceQuery{
		Resource:   resource,
		SharedWith: sharedWith,
	}
}
