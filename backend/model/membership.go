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

func (m *Membership) GetGroupKey() GroupKey {
	return NewGroupKey(m.GroupID)
}
