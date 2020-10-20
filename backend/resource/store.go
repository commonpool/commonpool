package resource

import "github.com/commonpool/backend/model"

type Query struct {
	Type      *model.ResourceType
	Query     *string
	Skip      int
	Take      int
	CreatedBy string
}

type QueryResult struct {
	Items      []model.Resource
	TotalCount int
	Skip       int
	Take       int
}

type Store interface {
	GetByKey(key model.ResourceKey, r *model.Resource) error
	Search(query Query) (*QueryResult, error)
	Delete(key model.ResourceKey) error
	Create(resource *model.Resource) error
	Update(resource *model.Resource) error
}
