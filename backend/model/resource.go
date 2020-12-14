package model

import (
	"github.com/commonpool/backend/pkg/exceptions"
	uuid "github.com/satori/go.uuid"
)

type ResourceKey struct {
	ID uuid.UUID
}

func NewResourceKey(id uuid.UUID) ResourceKey {
	return ResourceKey{
		ID: id,
	}
}

func ParseResourceKey(key string) (ResourceKey, error) {
	resourceUuid, err := uuid.FromString(key)
	if err != nil {
		response := exceptions.ErrInvalidResourceKey(key)
		return ResourceKey{}, &response
	}
	resourceKey := ResourceKey{
		ID: resourceUuid,
	}
	return resourceKey, nil
}

func (r *ResourceKey) GetUUID() uuid.UUID {
	return r.ID
}

func (r *ResourceKey) String() string {
	return r.ID.String()
}

type ResourceKeys struct {
	Items []ResourceKey
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
