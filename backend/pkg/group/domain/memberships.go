package domain

import (
	"github.com/commonpool/backend/pkg/keys"
)

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

func (m *Memberships) GetMembership(key keys.MembershipKey) (*Membership, bool) {
	for _, item := range m.Items {
		if item.GetKey() == key {
			return item, true
		}
	}
	return nil, false
}

func (m *Memberships) RemoveMembership(key keys.MembershipKey) {
	var idx *int
	for i, item := range m.Items {
		if item.GetKey() == key {
			idx = &i
			break
		}
	}
	if idx != nil {
		i := *idx
		copy(m.Items[i:], m.Items[i+1:])
		m.Items[len(m.Items)-1] = nil
		m.Items = m.Items[:len(m.Items)-1]
	}
}
