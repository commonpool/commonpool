package resource

import (
	ctx "context"
	"github.com/commonpool/backend/model"
)

type Store interface {
	GetByKey(ctx ctx.Context, getResourceByKeyQuery *GetResourceByKeyQuery) (*GetResourceByKeyResponse, error)
	GetByKeys(ctx ctx.Context, resourceKeys *model.ResourceKeys) (*GetResourceByKeysResponse, error)
	Search(searchResourcesQuery *SearchResourcesQuery) *SearchResourcesResponse
	Delete(deleteResourceQuery *DeleteResourceQuery) *DeleteResourceResponse
	Create(createResourceQuery *CreateResourceQuery) *CreateResourceResponse
	Update(updateResourceQuery *UpdateResourceQuery) *UpdateResourceResponse
}

type GetResourceByKeysQuery struct {
	ResourceKeys []model.ResourceKey
}

type SearchResourcesQuery struct {
	Type            *Type
	Query           *string
	Skip            int
	Take            int
	CreatedBy       string
	SharedWithGroup *model.GroupKey
}

func NewSearchResourcesQuery(query *string, resourceType *Type, skip int, take int, createdBy string, sharedWithGroup *model.GroupKey) *SearchResourcesQuery {
	return &SearchResourcesQuery{
		Type:            resourceType,
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
	Error      error
}

type GetResourceByKeyQuery struct {
	ResourceKey model.ResourceKey
}

func NewGetResourceByKeyQuery(resourceKey model.ResourceKey) *GetResourceByKeyQuery {
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

type DeleteResourceQuery struct {
	ResourceKey model.ResourceKey
}

type DeleteResourceResponse struct {
	Error error
}

type CreateResourceQuery struct {
	Resource   *Resource
	SharedWith *model.GroupKeys
}

func NewCreateResourceQuery(resource *Resource, sharedWith *model.GroupKeys) *CreateResourceQuery {
	return &CreateResourceQuery{
		Resource:   resource,
		SharedWith: sharedWith,
	}
}

type CreateResourceResponse struct {
	Error error
}

type UpdateResourceQuery struct {
	Resource   *Resource
	SharedWith *model.GroupKeys
}

func NewUpdateResourceQuery(resource *Resource, sharedWith *model.GroupKeys) *UpdateResourceQuery {
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
