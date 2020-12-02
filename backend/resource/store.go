package resource

import (
	ctx "context"
	"github.com/commonpool/backend/model"
)

type GetResourceByKeysResponse struct {
	Items *Resources
}

type GetResourceByKeysQuery struct {
	ResourceKeys []model.ResourceKey
}

func NewGetResourceByKeysQuery(keys []model.ResourceKey) *GetResourceByKeysQuery {
	return &GetResourceByKeysQuery{
		ResourceKeys: keys,
	}
}

type Store interface {
	GetByKey(ctx ctx.Context, getResourceByKeyQuery *GetResourceByKeyQuery) *GetResourceByKeyResponse
	GetByKeys(getResourceByKeysQuery *GetResourceByKeysQuery) (*GetResourceByKeysResponse, error)
	Search(searchResourcesQuery *SearchResourcesQuery) *SearchResourcesResponse
	Delete(deleteResourceQuery *DeleteResourceQuery) *DeleteResourceResponse
	Create(createResourceQuery *CreateResourceQuery) *CreateResourceResponse
	Update(updateResourceQuery *UpdateResourceQuery) *UpdateResourceResponse
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

func NewSearchResourcesResponseSuccess(resources *Resources, sharings *Sharings, totalCount int, skip int, take int) *SearchResourcesResponse {
	return &SearchResourcesResponse{
		Resources:  resources,
		Sharings:   sharings,
		TotalCount: totalCount,
		Skip:       skip,
		Take:       take,
	}
}

func NewSearchResourcesResponseError(err error) *SearchResourcesResponse {
	return &SearchResourcesResponse{
		Resources:  nil,
		Sharings:   nil,
		TotalCount: -1,
		Skip:       -1,
		Take:       -1,
		Error:      err,
	}
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
	Error    error
}

func NewGetResourceByKeyResponseError(err error) *GetResourceByKeyResponse {
	return &GetResourceByKeyResponse{
		Error: err,
	}
}

func NewGetResourceByKeyResponseSuccess(resource *Resource, sharings *Sharings) *GetResourceByKeyResponse {
	return &GetResourceByKeyResponse{
		Resource: resource,
		Sharings: sharings,
	}
}

type DeleteResourceQuery struct {
	ResourceKey model.ResourceKey
}

func NewDeleteResourceQuery(key model.ResourceKey) *DeleteResourceQuery {
	return &DeleteResourceQuery{
		ResourceKey: key,
	}
}

type DeleteResourceResponse struct {
	Error error
}

func NewDeleteResourceResponse(err error) *DeleteResourceResponse {
	return &DeleteResourceResponse{
		Error: err,
	}
}

type CreateResourceQuery struct {
	Resource   *Resource
	SharedWith []model.GroupKey
}

func NewCreateResourceQuery(resource *Resource, sharedWith []model.GroupKey) *CreateResourceQuery {
	return &CreateResourceQuery{
		Resource:   resource,
		SharedWith: sharedWith,
	}
}

type CreateResourceResponse struct {
	Error error
}

func NewCreateResourceResponse(err error) *CreateResourceResponse {
	return &CreateResourceResponse{
		Error: err,
	}
}

type UpdateResourceQuery struct {
	Resource   *Resource
	SharedWith []model.GroupKey
}

func NewUpdateResourceQuery(resource *Resource, sharedWith []model.GroupKey) *UpdateResourceQuery {
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
