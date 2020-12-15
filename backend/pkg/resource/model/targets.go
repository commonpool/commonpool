package model

import (
	model2 "github.com/commonpool/backend/pkg/group/model"
	"github.com/commonpool/backend/pkg/user/model"
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

func (t *Targets) GetUserKeys() *model.UserKeys {
	var userKeys []model.UserKey
	for _, target := range t.Items {
		if !target.IsForUser() {
			continue
		}
		userKeys = append(userKeys, target.GetUserKey())
	}
	return model.NewUserKeys(userKeys)
}

func (t *Targets) GetGroupKeys() *model2.GroupKeys {
	var groupKeys []model2.GroupKey
	for _, target := range t.Items {
		if !target.IsForGroup() {
			continue
		}
		groupKeys = append(groupKeys, target.GetGroupKey())
	}
	return model2.NewGroupKeys(groupKeys)
}
