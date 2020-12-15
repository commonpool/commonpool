package model

import (
	"github.com/commonpool/backend/pkg/exceptions"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	usermodel "github.com/commonpool/backend/pkg/user/model"
)

type Target struct {
	UserKey  *usermodel.UserKey
	GroupKey *groupmodel.GroupKey
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

func (t Target) GetGroupKey() groupmodel.GroupKey {
	return *t.GroupKey
}
func (t Target) GetUserKey() usermodel.UserKey {
	return *t.UserKey
}

func NewUserTarget(userKey usermodel.UserKey) *Target {
	return &Target{
		UserKey: &userKey,
		Type:    UserTarget,
	}
}
func NewGroupTarget(groupKey groupmodel.GroupKey) *Target {
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
		return "", exceptions.ErrInvalidTargetType
	}
}
