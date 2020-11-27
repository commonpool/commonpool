package model

import (
	errs "github.com/commonpool/backend/errors"
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

func ParseResourceKey(key string) (*ResourceKey, error) {
	resourceUuid, err := uuid.FromString(key)
	if err != nil {
		response := errs.ErrInvalidResourceKey(key)
		return nil, &response
	}
	resourceKey := ResourceKey{
		ID: resourceUuid,
	}
	return &resourceKey, nil
}

func (r *ResourceKey) GetUUID() uuid.UUID {
	return r.ID
}

func (r *ResourceKey) String() string {
	return r.ID.String()
}
