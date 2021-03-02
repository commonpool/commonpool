package keys

import (
	"encoding/json"
	"github.com/satori/go.uuid"
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
func GenerateOfferItemKey() OfferItemKey {
	return OfferItemKey{ID: uuid.NewV4()}
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

func (k OfferItemKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.ID.String())
}

func (k *OfferItemKey) UnmarshalJSON(data []byte) error {
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
