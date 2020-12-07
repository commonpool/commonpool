package model

import (
	uuid "github.com/satori/go.uuid"
)

type OfferItemKey struct {
	ID uuid.UUID
}

func (o OfferItemKey) String() string {
	return o.ID.String()
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

type OfferItemKeys struct {
	Items []OfferItemKey
}

func (t *OfferItemKeys) Strings() []string {
	var strings []string
	for _, item := range t.Items {
		strings = append(strings, item.String())
	}
	return strings
}

func NewOfferItemKeys(items []OfferItemKey) *OfferItemKeys {
	copied := make([]OfferItemKey, len(items))
	copy(copied, items)
	return &OfferItemKeys{
		Items: copied,
	}
}
