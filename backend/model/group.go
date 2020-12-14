package model

import (
	"github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/utils"
	uuid "github.com/satori/go.uuid"
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

func (k GroupKey) GetChannelKey() ChannelKey {
	shortUid := utils.ShortUuid(k.ID)
	return ChannelKey{
		ID: shortUid,
	}
}

func NewGroupKey(id uuid.UUID) GroupKey {
	return GroupKey{ID: id}
}

func (k GroupKey) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("group_id", k.ID.String())
	return nil
}

type GroupKeys struct {
	Items []GroupKey
}

func NewGroupKeys(groupKeys []GroupKey) *GroupKeys {
	copied := make([]GroupKey, len(groupKeys))
	copy(copied, groupKeys)
	return &GroupKeys{
		Items: copied,
	}
}

func NewEmptyGroupKeys() *GroupKeys {
	return NewGroupKeys([]GroupKey{})
}

func (k GroupKeys) Strings() []string {
	var groupKeys []string
	for _, groupKey := range k.Items {
		groupKeys = append(groupKeys, groupKey.String())
	}
	if groupKeys == nil {
		groupKeys = []string{}
	}
	return groupKeys
}

func (k GroupKeys) Contains(groupKey GroupKey) bool {
	for _, gk := range k.Items {
		if groupKey == gk {
			return true
		}
	}
	return false
}

func (i *GroupKeys) IsEmpty() bool {
	return i.Items == nil || len(i.Items) == 0
}

func (i *GroupKeys) Count() int {
	if i.Items == nil {
		return 0
	}
	return len(i.Items)
}

func ParseGroupKey(value string) (GroupKey, error) {
	offerId, err := uuid.FromString(value)
	if err != nil {
		return GroupKey{}, errors.ErrInvalidGroupId
	}
	return NewGroupKey(offerId), err
}
