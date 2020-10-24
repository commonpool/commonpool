package model

import (
	errs "github.com/commonpool/backend/errors"
	uuid "github.com/satori/go.uuid"
)

type ResourceKey struct {
	uuid uuid.UUID
}

func NewResourceKey() ResourceKey {
	return ResourceKey{
		uuid: uuid.NewV4(),
	}
}

func ParseResourceKey(key string) (*ResourceKey, error) {
	resourceUuid, err := uuid.FromString(key)
	if err != nil {
		response := errs.ErrInvalidResourceKey(key)
		return nil, &response
	}
	resourceKey := ResourceKey{
		uuid: resourceUuid,
	}
	return &resourceKey, nil
}

func (r *ResourceKey) GetUUID() uuid.UUID {
	return r.uuid
}

func (r *ResourceKey) String() string {
	return r.uuid.String()
}
