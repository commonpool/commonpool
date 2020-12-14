package resource

import (
	"fmt"
	"github.com/commonpool/backend/model"
)

type Resources struct {
	ItemMap map[model.ResourceKey]*Resource
	Items   []*Resource
}

func NewResources(items []*Resource) *Resources {
	rsMap := map[model.ResourceKey]*Resource{}
	for _, item := range items {
		rsMap[item.GetKey()] = item
	}
	return &Resources{
		Items:   items,
		ItemMap: rsMap,
	}
}

func NewEmptyResources() *Resources {
	return &Resources{
		ItemMap: map[model.ResourceKey]*Resource{},
		Items:   []*Resource{},
	}
}

func (r *Resources) GetResource(key model.ResourceKey) (*Resource, error) {
	rs, ok := r.ItemMap[key]
	if !ok {
		return nil, fmt.Errorf("resource not found")
	}
	return rs, nil
}

func (r *Resources) Append(resource *Resource) *Resources {
	items := append(r.Items, resource)
	return NewResources(items)
}

func (r *Resources) Contains(resource *Resource) bool {
	return r.ContainsKey(resource.GetKey())
}

func (r *Resources) ContainsKey(key model.ResourceKey) bool {
	_, ok := r.ItemMap[key]
	return ok
}

func (r *Resources) GetKeys() *model.ResourceKeys {
	var resourceKeys []model.ResourceKey
	for _, resource := range r.Items {
		resourceKeys = append(resourceKeys, resource.GetKey())
	}
	if resourceKeys == nil {
		resourceKeys = []model.ResourceKey{}
	}
	return model.NewResourceKeys(resourceKeys)
}
