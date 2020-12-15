package resource

import (
	"context"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
)

type Store interface {
	GetByKey(ctx context.Context, getResourceByKeyQuery *GetResourceByKeyQuery) (*GetResourceByKeyResponse, error)
	GetByKeys(ctx context.Context, resourceKeys *resourcemodel.ResourceKeys) (*GetResourceByKeysResponse, error)
	Search(ctx context.Context, searchResourcesQuery *SearchResourcesQuery) *SearchResourcesResponse
	Delete(deleteResourceQuery *DeleteResourceQuery) *DeleteResourceResponse
	Create(createResourceQuery *CreateResourceQuery) *CreateResourceResponse
	Update(updateResourceQuery *UpdateResourceQuery) *UpdateResourceResponse
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
	Error      error
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

type DeleteResourceQuery struct {
	ResourceKey resourcemodel.ResourceKey
}

type DeleteResourceResponse struct {
	Error error
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

type CreateResourceResponse struct {
	Error error
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

type UpdateResourceResponse struct {
	Error error
}

func NewUpdateResourceResponse(err error) *UpdateResourceResponse {
	return &UpdateResourceResponse{
		Error: err,
	}
}
