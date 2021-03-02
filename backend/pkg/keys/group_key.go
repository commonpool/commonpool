package keys

import (
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/satori/go.uuid"
	"go.uber.org/zap/zapcore"
)

type GroupKey struct {
	ID uuid.UUID
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
func GenerateGroupKey() GroupKey {
	return GroupKey{ID: uuid.NewV4()}
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
