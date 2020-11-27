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
