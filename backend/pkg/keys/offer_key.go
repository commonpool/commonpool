package keys

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"go.uber.org/zap/zapcore"
)

type OfferKey struct {
	ID uuid.UUID
}

func NewOfferKey(id uuid.UUID) OfferKey {
	return OfferKey{ID: id}
}

func GenerateOfferKey() OfferKey {
	return OfferKey{ID: uuid.NewV4()}
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

func (ok OfferKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(ok.ID.String())
}

func (k *OfferKey) UnmarshalJSON(data []byte) error {
	var uid string
	if err := json.Unmarshal(data, &uid); err != nil {
		return err
	}
	id, err := uuid.FromString(uid)
	if err != nil {
		return err
	}
	k.ID = id
	return nil
}
