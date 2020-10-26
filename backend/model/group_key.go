package model

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
)

type GroupKey struct {
	ID uuid.UUID
}

func NewGroupKey(id uuid.UUID) GroupKey {
	return GroupKey{ID: id}
}

func ParseGroupKey(value string) (GroupKey, error) {
	offerId, err := uuid.FromString(value)
	if err != nil {
		return GroupKey{}, fmt.Errorf("cannot parse group key: %s", err.Error())
	}
	return NewGroupKey(offerId), err
}
