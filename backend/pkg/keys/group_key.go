package keys

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/satori/go.uuid"
	"go.uber.org/zap/zapcore"
)

type GroupKey struct {
	ID uuid.UUID
}

func (g GroupKey) GetGroupKey() GroupKey {
	return g
}

type GroupKeyGetter interface {
	GetGroupKey() GroupKey
}

var _ zapcore.ObjectMarshaler = &GroupKey{}

func (k GroupKey) String() string {
	return k.ID.String()
}

func (k GroupKey) Equals(g GroupKey) bool {
	return k.ID == g.ID
}

func NewGroupKey(id uuid.UUID) GroupKey {
	return GroupKey{ID: id}
}

func (k GroupKey) Target() *Target {
	return NewGroupTarget(k)
}

func GenerateGroupKey() GroupKey {
	return GroupKey{ID: uuid.NewV4()}
}

func (k GroupKey) StreamKey() StreamKey {
	return NewStreamKey("group", k.String())
}

func (k GroupKey) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("group_id", k.ID.String())
	return nil
}

func ParseGroupKey(value string) (GroupKey, error) {
	offerId, err := uuid.FromString(value)
	if err != nil {
		return GroupKey{}, exceptions.ErrInvalidGroupId
	}
	return NewGroupKey(offerId), err
}

func (k GroupKey) GetFrontendLink() string {
	return fmt.Sprintf("<commonpool-group id='%s'><commonpool-group>", k.String())
}

func (k GroupKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.ID.String())
}

func (k *GroupKey) UnmarshalJSON(data []byte) error {
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

func (k *GroupKey) Scan(value interface{}) error {
	keyValue, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal string value:", value))
	}
	uid, err := uuid.FromString(keyValue)
	if err != nil {
		return err
	}
	*k = NewGroupKey(uid)
	return nil
}

func (k GroupKey) Value() (driver.Value, error) {
	return driver.String.ConvertValue(k.ID.String())
}

func (k GroupKey) GormDataType() string {
	return "varchar(128)"
}
