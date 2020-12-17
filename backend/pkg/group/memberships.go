package group

import "github.com/commonpool/backend/pkg/keys"

type Memberships struct {
	Items []*Membership
}

func NewMemberships(items []*Membership) *Memberships {
	return &Memberships{Items: items}
}

func (m *Memberships) ContainsMembershipForGroup(groupKey keys.GroupKey) bool {
	for _, item := range m.Items {
		if item.GetGroupKey().Equals(groupKey) {
			return true
		}
	}
	return false
}
