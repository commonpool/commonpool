package group

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

func (gk GroupKeys) Strings() []string {
	var groupKeys []string
	for _, groupKey := range gk.Items {
		groupKeys = append(groupKeys, groupKey.String())
	}
	if groupKeys == nil {
		groupKeys = []string{}
	}
	return groupKeys
}

func (gk GroupKeys) Contains(groupKey GroupKey) bool {
	for _, gk := range gk.Items {
		if groupKey == gk {
			return true
		}
	}
	return false
}

func (gk *GroupKeys) IsEmpty() bool {
	return gk.Items == nil || len(gk.Items) == 0
}

func (gk *GroupKeys) Count() int {
	if gk.Items == nil {
		return 0
	}
	return len(gk.Items)
}
