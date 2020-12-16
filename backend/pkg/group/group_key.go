package group

import (
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
