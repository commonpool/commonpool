package model

import (
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/user/usermodel"
)

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

func (t *Targets) GetUserKeys() *usermodel.UserKeys {
	var userKeys []usermodel.UserKey
	for _, target := range t.Items {
		if !target.IsForUser() {
			continue
		}
		userKeys = append(userKeys, target.GetUserKey())
	}
	return usermodel.NewUserKeys(userKeys)
}

func (t *Targets) GetGroupKeys() *group.GroupKeys {
	var groupKeys []group.GroupKey
	for _, target := range t.Items {
		if !target.IsForGroup() {
			continue
		}
		groupKeys = append(groupKeys, target.GetGroupKey())
	}
	return group.NewGroupKeys(groupKeys)
}
