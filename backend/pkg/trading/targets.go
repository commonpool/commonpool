package trading

import (
	"github.com/commonpool/backend/pkg/keys"
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

func (t *Targets) GetUserKeys() *keys.UserKeys {
	var userKeys []keys.UserKey
	for _, target := range t.Items {
		if !target.IsForUser() {
			continue
		}
		userKeys = append(userKeys, target.GetUserKey())
	}
	return keys.NewUserKeys(userKeys)
}

func (t *Targets) GetGroupKeys() *keys.GroupKeys {
	var groupKeys []keys.GroupKey
	for _, target := range t.Items {
		if !target.IsForGroup() {
			continue
		}
		groupKeys = append(groupKeys, target.GetGroupKey())
	}
	return keys.NewGroupKeys(groupKeys)
}
