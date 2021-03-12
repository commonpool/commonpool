package keys

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
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

func (ok *OfferItemKey) UnmarshalParam(param string) error {
	offerItemID, err := uuid.FromString(param)
	if err != nil {
		return err
	}
	ok.ID = offerItemID
	return nil
}

func (k *OfferItemKey) Scan(value interface{}) error {
	keyValue, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal string value:", value))
	}
	uid, err := uuid.FromString(keyValue)
	if err != nil {
		return err
	}
	*k = NewOfferItemKey(uid)
	return nil
}

func (k OfferItemKey) Value() (driver.Value, error) {
	return driver.String.ConvertValue(k.String())
}

func (k OfferItemKey) GormDataType() string {
	return "varchar(128)"
}
