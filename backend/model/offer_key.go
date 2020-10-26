package model

import uuid "github.com/satori/go.uuid"

type OfferKey struct {
	ID uuid.UUID
}

func NewOfferKey(id uuid.UUID) OfferKey {
	return OfferKey{ID: id}
}

func ParseOfferKey(value string) (OfferKey, error) {
	offerId, err := uuid.FromString(value)
	if err != nil {
		return OfferKey{}, err
	}
	return NewOfferKey(offerId), err
}
