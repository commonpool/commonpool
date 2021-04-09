package api

import "time"

type Membership struct {
	GroupID         string `gorm:"primaryKey"`
	UserID          string `gorm:"primaryKey"`
	Permission      MembershipPermission
	User            *User
	Group           *Group
	MemberConfirmed bool
	GroupConfirmed  bool
	CreatedAt       time.Time
}

func (m *Membership) IsAdmin() bool {
	return m.IsActive() && (m.Permission == Owner || m.Permission == "admin")
}

func (m *Membership) IsOwner() bool {
	return m.IsActive() && m.Permission == Owner
}

func (m *Membership) IsActive() bool {
	return m.GroupConfirmed && m.MemberConfirmed
}

type MembershipPermission string

func (m MembershipPermission) Gte(o MembershipPermission) bool {
	if m == Owner {
		return true
	}
	if m == Admin && (o == Admin || o == Member || o == None) {
		return true
	}
	if m == Member && (o == Member || o == None) {
		return true
	}
	if m == None && o == None {
		return true
	}
	return false
}

const (
	None   MembershipPermission = "none"
	Member MembershipPermission = "member"
	Admin  MembershipPermission = "admin"
	Owner  MembershipPermission = "owner"
)
