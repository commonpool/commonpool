package model

import (
	uuid "github.com/satori/go.uuid"
)

type OfferItemKey struct {
	ID uuid.UUID
}

func NewOfferItemKey(id uuid.UUID) OfferItemKey {
	return OfferItemKey{ID: id}
}

func ParseOfferItemKey(str string) (OfferItemKey, error) {
	uid, err := uuid.FromString(str)
	if err != nil {
		return OfferItemKey{}, err
	}
	return NewOfferItemKey(uid), nil
}

func MustParseOfferItemKey(str string) OfferItemKey {
	uid, err := uuid.FromString(str)
	if err != nil {
		panic(err)
	}
	return NewOfferItemKey(uid)
}
