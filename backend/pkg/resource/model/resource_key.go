package model

import (
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/satori/go.uuid"
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
