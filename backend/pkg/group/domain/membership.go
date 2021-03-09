package domain

import (
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type PermissionLevel int

const (
	None PermissionLevel = iota
	Member
	Admin
	Owner
)

func (p PermissionLevel) IsMember() bool {
	return p >= Member
}

func (p PermissionLevel) IsAdmin() bool {
	return p >= Admin
}

func (p PermissionLevel) IsOwner() bool {
	return p >= Owner
}

type Membership struct {
	PermissionLevel
	Key       keys.MembershipKey
	Status    MembershipStatus
	CreatedAt time.Time
}

func (m *Membership) HasGroupConfirmed() bool {
	return m.Status == PendingUserMembershipStatus || m.Status == ApprovedMembershipStatus
}

func (m *Membership) HasUserConfirmed() bool {
	return m.Status == PendingGroupMembershipStatus || m.Status == ApprovedMembershipStatus
}

func (m *Membership) HasBothPartiesConfirmed() bool {
	return m.Status == ApprovedMembershipStatus
}

func (m *Membership) GetGroupKey() keys.GroupKey {
	return m.Key.GroupKey
}

func (m *Membership) GetUserKey() keys.UserKey {
	return m.Key.UserKey
}

func (m *Membership) GetKey() keys.MembershipKey {
	return keys.NewMembershipKey(m, m.GetUserKey())
}

func (m *Membership) GetStatus() MembershipStatus {
	return m.Status
}
