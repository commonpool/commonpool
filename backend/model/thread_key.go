package model

import uuid "github.com/satori/go.uuid"

type ThreadKey struct {
	uuid uuid.UUID
}

func NewThreadKey() ThreadKey {
	return ThreadKey{
		uuid: uuid.NewV4(),
	}
}

func ParseThreadKey(key string) (*ThreadKey, error) {
	resourceUuid, err := uuid.FromString(key)
	if err != nil {
		return nil, err
	}
	threadKey := ThreadKey{
		uuid: resourceUuid,
	}
	return &threadKey, nil
}

func (r *ThreadKey) GetUUID() uuid.UUID {
	return r.uuid
}

func (r *ThreadKey) String() string {
	return r.uuid.String()
}
