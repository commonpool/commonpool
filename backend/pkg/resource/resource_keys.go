package resource

import "github.com/commonpool/backend/pkg/keys"

type ResourceKeys struct {
	Items []keys.ResourceKey
}

func (k ResourceKeys) Count() int {
	return len(k.Items)
}

func (k ResourceKeys) IsEmpty() bool {
	return k.Count() == 0
}

func (k ResourceKeys) Strings() []string {
	var strings []string
	for _, item := range k.Items {
		strings = append(strings, item.String())
	}
	if strings == nil {
		strings = []string{}
	}
	return strings
}

func NewResourceKeys(resourceKeys []keys.ResourceKey) *ResourceKeys {
	copied := make([]keys.ResourceKey, len(resourceKeys))
	copy(copied, resourceKeys)
	return &ResourceKeys{
		Items: copied,
	}
}

func NewEmptyResourceKeys() *ResourceKeys {
	return &ResourceKeys{
		Items: []keys.ResourceKey{},
	}
}
