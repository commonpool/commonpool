package model

import uuid "github.com/satori/go.uuid"

type OfferItemKey struct {
	ID      uuid.UUID
	OfferID uuid.UUID
}

func NewOfferItemKey(id uuid.UUID, offerKey OfferKey) OfferItemKey {
	return OfferItemKey{ID: id, OfferID: offerKey.ID}
}
