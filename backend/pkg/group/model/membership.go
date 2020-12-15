package model

import (
	usermodel "github.com/commonpool/backend/pkg/user/model"
	"time"
)

type Membership struct {
	Key            MembershipKey
	IsMember       bool
	IsAdmin        bool
	IsOwner        bool
	GroupConfirmed bool
	UserConfirmed  bool
	CreatedAt      time.Time
	IsDeactivated  bool
}

func (m *Membership) GetGroupKey() GroupKey {
	return m.Key.GroupKey
}

func (m *Membership) GetUserKey() usermodel.UserKey {
	return m.Key.UserKey
}

func (m *Membership) GetKey() MembershipKey {
	return NewMembershipKey(m.GetGroupKey(), m.GetUserKey())
}
