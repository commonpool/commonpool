package group

import (
	"github.com/commonpool/backend/model"
	"time"
)

type Membership struct {
	Key            model.MembershipKey
	IsMember       bool
	IsAdmin        bool
	IsOwner        bool
	GroupConfirmed bool
	UserConfirmed  bool
	CreatedAt      time.Time
	IsDeactivated  bool
}

func (m *Membership) GetGroupKey() model.GroupKey {
	return m.Key.GroupKey
}

func (m *Membership) GetUserKey() model.UserKey {
	return m.Key.UserKey
}

func (m *Membership) GetKey() model.MembershipKey {
	return model.NewMembershipKey(m.GetGroupKey(), m.GetUserKey())
}
