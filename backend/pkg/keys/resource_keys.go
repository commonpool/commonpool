package keys

type ResourceKeys struct {
	Items []ResourceKey
}

func (k ResourceKeys) Count() int {
	return len(k.Items)
}

func (k ResourceKeys) IsEmpty() bool {
	return k.Count() == 0
}

func (k *ResourceKeys) Append(resourceKey ResourceKey) *ResourceKeys {
	newResourceKeys := make([]ResourceKey, len(k.Items)+1)
	copy(newResourceKeys, k.Items)
	newResourceKeys[len(k.Items)] = resourceKey
	return NewResourceKeys(newResourceKeys)
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

func NewResourceKeys(resourceKeys []ResourceKey) *ResourceKeys {
	copied := make([]ResourceKey, len(resourceKeys))
	copy(copied, resourceKeys)
	return &ResourceKeys{
		Items: copied,
	}
}

func NewEmptyResourceKeys() *ResourceKeys {
	return &ResourceKeys{
		Items: []ResourceKey{},
	}
}
