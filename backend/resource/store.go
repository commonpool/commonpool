package resource

import "github.com/commonpool/backend/model"

type Store interface {
	GetByKey(getResourceByKeyQuery *GetResourceByKeyQuery) *GetResourceByKeyResponse
	Search(searchResourcesQuery *SearchResourcesQuery) *SearchResourcesResponse
	Delete(deleteResourceQuery *DeleteResourceQuery) *DeleteResourceResponse
	Create(createResourceQuery *CreateResourceQuery) *CreateResourceResponse
	Update(updateResourceQuery *UpdateResourceQuery) *UpdateResourceResponse
}

type SearchResourcesQuery struct {
	Type            *model.ResourceType
	Query           *string
	Skip            int
	Take            int
	CreatedBy       string
	SharedWithGroup *model.GroupKey
}

func NewSearchResourcesQuery(query *string, resourceType *model.ResourceType, skip int, take int, createdBy string, sharedWithGroup *model.GroupKey) *SearchResourcesQuery {
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
	Resources  *model.Resources
	Sharings   *model.ResourceSharings
	TotalCount int
	Skip       int
	Take       int
	Error      error
}

func NewSearchResourcesResponseSuccess(resources *model.Resources, sharings *model.ResourceSharings, totalCount int, skip int, take int) *SearchResourcesResponse {
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
	Resource *model.Resource
	Sharings *model.ResourceSharings
	Error    error
}

func NewGetResourceByKeyResponseError(err error) *GetResourceByKeyResponse {
	return &GetResourceByKeyResponse{
		Error: err,
	}
}

func NewGetResourceByKeyResponseSuccess(resource *model.Resource, sharings *model.ResourceSharings) *GetResourceByKeyResponse {
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
	Resource   *model.Resource
	SharedWith []model.GroupKey
}

func NewCreateResourceQuery(resource *model.Resource) *CreateResourceQuery {
	return &CreateResourceQuery{
		Resource: resource,
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
	Resource   *model.Resource
	SharedWith []model.GroupKey
}

func NewUpdateResourceQuery(resource *model.Resource, sharedWith []model.GroupKey) *UpdateResourceQuery {
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
