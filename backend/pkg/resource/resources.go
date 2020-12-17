package resource

import (
	"fmt"
	"github.com/commonpool/backend/pkg/keys"
)

type Resources struct {
	ItemMap map[keys.ResourceKey]*Resource
	Items   []*Resource
}

func NewResources(items []*Resource) *Resources {
	rsMap := map[keys.ResourceKey]*Resource{}
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
		ItemMap: map[keys.ResourceKey]*Resource{},
		Items:   []*Resource{},
	}
}

func (r *Resources) GetResource(key keys.ResourceKey) (*Resource, error) {
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

func (r *Resources) ContainsKey(key keys.ResourceKey) bool {
	_, ok := r.ItemMap[key]
	return ok
}

func (r *Resources) GetKeys() *ResourceKeys {
	var resourceKeys []keys.ResourceKey
	for _, resource := range r.Items {
		resourceKeys = append(resourceKeys, resource.GetKey())
	}
	if resourceKeys == nil {
		resourceKeys = []keys.ResourceKey{}
	}
	return NewResourceKeys(resourceKeys)
}
