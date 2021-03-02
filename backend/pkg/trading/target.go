package trading

import (
	"github.com/commonpool/backend/pkg/keys"
)

type Target struct {
	UserKey  *keys.UserKey  `json:"userId,omitempty"`
	GroupKey *keys.GroupKey `json:"groupId,omitempty"`
	Type     TargetType     `json:"type"`
}

func (t Target) Equals(target *Target) bool {

	if t.Type != target.Type {
		return false
	}

	if t.Type == GroupTarget {
		return *t.GroupKey == *target.GroupKey
	}

	return *t.UserKey == *target.UserKey
}

func (t Target) IsForGroup() bool {
	return t.Type == GroupTarget
}

func (t Target) IsForUser() bool {
	return t.Type == UserTarget
}

func (t Target) GetGroupKey() keys.GroupKey {
	return *t.GroupKey
}
func (t Target) GetUserKey() keys.UserKey {
	return *t.UserKey
}

func NewUserTarget(userKey keys.UserKey) *Target {
	return &Target{
		UserKey: &userKey,
		Type:    UserTarget,
	}
}
func NewGroupTarget(groupKey keys.GroupKey) *Target {
	return &Target{
		GroupKey: &groupKey,
		Type:     GroupTarget,
	}
}
