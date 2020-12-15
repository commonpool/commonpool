package model

import (
	"fmt"
)

type Resources struct {
	ItemMap map[ResourceKey]*Resource
	Items   []*Resource
}

func NewResources(items []*Resource) *Resources {
	rsMap := map[ResourceKey]*Resource{}
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
		ItemMap: map[ResourceKey]*Resource{},
		Items:   []*Resource{},
	}
}

func (r *Resources) GetResource(key ResourceKey) (*Resource, error) {
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

func (r *Resources) ContainsKey(key ResourceKey) bool {
	_, ok := r.ItemMap[key]
	return ok
}

func (r *Resources) GetKeys() *ResourceKeys {
	var resourceKeys []ResourceKey
	for _, resource := range r.Items {
		resourceKeys = append(resourceKeys, resource.GetKey())
	}
	if resourceKeys == nil {
		resourceKeys = []ResourceKey{}
	}
	return NewResourceKeys(resourceKeys)
}
