package model

import (
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap/zapcore"
)

type OfferKey struct {
	ID uuid.UUID
}

func NewOfferKey(id uuid.UUID) OfferKey {
	return OfferKey{ID: id}
}

func (ok OfferKey) String() string {
	return ok.ID.String()
}

func ParseOfferKey(value string) (OfferKey, error) {
	offerId, err := uuid.FromString(value)
	if err != nil {
		return OfferKey{}, err
	}
	return NewOfferKey(offerId), err
}

//goland:noinspection GoUnusedExportedFunction
func MustParseOfferKey(value string) OfferKey {
	offerId, err := uuid.FromString(value)
	if err != nil {
		panic(err)
	}
	return NewOfferKey(offerId)
}

func (ok OfferKey) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("offer_id", ok.String())
	return nil
}

type OfferKeys struct {
	Items []OfferKey
}

func NewOfferKeys(items []OfferKey) *OfferKeys {
	copied := make([]OfferKey, len(items))
	copy(copied, items)
	return &OfferKeys{
		Items: copied,
	}
}

func (o *OfferKeys) Strings() []string {
	var result []string
	for _, item := range o.Items {
		result = append(result, item.String())
	}
	if result == nil {
		result = []string{}
	}
	return result
}
