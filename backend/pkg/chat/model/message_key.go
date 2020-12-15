package model

import "github.com/satori/go.uuid"

type MessageKey struct {
	uuid uuid.UUID
}

func NewMessageKey(uuid uuid.UUID) MessageKey {
	return MessageKey{
		uuid: uuid,
	}
}

func ParseMessageKey(key string) (*MessageKey, error) {
	resourceUuid, err := uuid.FromString(key)
	if err != nil {
		return nil, err
	}
	messageKey := MessageKey{
		uuid: resourceUuid,
	}
	return &messageKey, nil
}

func (r *MessageKey) GetUUID() uuid.UUID {
	return r.uuid
}

func (r *MessageKey) String() string {
	return r.uuid.String()
}
