package model

import uuid "github.com/satori/go.uuid"

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
		return nil, err
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
