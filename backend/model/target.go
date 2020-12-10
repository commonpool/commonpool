package model

import "github.com/commonpool/backend/errors"

type Target struct {
	UserKey  *UserKey
	GroupKey *GroupKey
	Type     TargetType
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

func (t Target) GetGroupKey() GroupKey {
	return *t.GroupKey
}
func (t Target) GetUserKey() UserKey {
	return *t.UserKey
}

func NewUserTarget(userKey UserKey) *Target {
	return &Target{
		UserKey: &userKey,
		Type:    UserTarget,
	}
}
func NewGroupTarget(groupKey GroupKey) *Target {
	return &Target{
		GroupKey: &groupKey,
		Type:     GroupTarget,
	}
}

type Targets struct {
	Items []*Target
}

func NewTargets(items []*Target) *Targets {
	copied := make([]*Target, len(items))
	copy(copied, items)
	return &Targets{
		Items: copied,
	}
}

func NewEmptyTargets() *Targets {
	return &Targets{
		Items: []*Target{},
	}
}

type TargetType string

const (
	UserTarget  TargetType = "user"
	GroupTarget TargetType = "group"
)

func (a TargetType) IsGroup() bool {
	return a == GroupTarget
}

func (a TargetType) IsUser() bool {
	return a == UserTarget
}

func ParseOfferItemTargetType(str string) (TargetType, error) {
	if str == "user" {
		return UserTarget, nil
	} else if str == "group" {
		return GroupTarget, nil
	} else {
		return "", errors.ErrInvalidTargetType
	}
}
