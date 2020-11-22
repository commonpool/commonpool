package model

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type PermissionType int

const (
	MemberPermission PermissionType = iota
	AdminPermission
)

type Membership struct {
	GroupID        uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID         string    `gorm:"primary_key"`
	IsMember       bool
	IsAdmin        bool
	IsOwner        bool
	GroupConfirmed bool
	UserConfirmed  bool
	CreatedAt      time.Time
	IsDeactivated  bool
}

func NewEmptyMembership(
	membershipKey MembershipKey,
) Membership {
	return Membership{
		GroupID:        membershipKey.GroupKey.ID,
		UserID:         membershipKey.UserKey.String(),
		IsMember:       false,
		IsAdmin:        false,
		IsOwner:        false,
		GroupConfirmed: false,
		UserConfirmed:  false,
		IsDeactivated:  false,
	}
}

func (m *Membership) GetGroupKey() GroupKey {
	return NewGroupKey(m.GroupID)
}

func (m *Membership) GetUserKey() UserKey {
	return NewUserKey(m.UserID)
}

func (m *Membership) GetKey() MembershipKey {
	return NewMembershipKey(m.GetGroupKey(), m.GetUserKey())
}

type Memberships struct {
	Items []Membership
}

func NewMemberships(items []Membership) Memberships {
	return Memberships{Items: items}
}

func (m *Memberships) ContainsMembershipForGroup(groupKey GroupKey) bool {
	for _, item := range m.Items {
		if item.GetGroupKey().Equals(groupKey) {
			return true
		}
	}
	return false
}
