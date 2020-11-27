package model

import (
	"fmt"
	"github.com/commonpool/backend/utils"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap/zapcore"
)

type GroupKey struct {
	ID uuid.UUID
}

var _ zapcore.ObjectMarshaler = &GroupKey{}

func (k GroupKey) Equals(g GroupKey) bool {
	return k.ID == g.ID
}

func (k GroupKey) GetChannelKey() ChannelKey {
	shortUid := utils.ShortUuid(k.ID)
	return ChannelKey{
		ID: shortUid,
	}
}

func NewGroupKey(id uuid.UUID) GroupKey {
	return GroupKey{ID: id}
}

func ParseGroupKey(value string) (GroupKey, error) {
	offerId, err := uuid.FromString(value)
	if err != nil {
		return GroupKey{}, fmt.Errorf("cannot parse group key: %s", err.Error())
	}
	return NewGroupKey(offerId), err
}

func (k GroupKey) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("group_id", k.ID.String())
	return nil
}
