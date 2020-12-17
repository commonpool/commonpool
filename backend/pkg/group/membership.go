package group

import (
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type Membership struct {
	Key            keys.MembershipKey
	IsMember       bool
	IsAdmin        bool
	IsOwner        bool
	GroupConfirmed bool
	UserConfirmed  bool
	CreatedAt      time.Time
	IsDeactivated  bool
}

func (m *Membership) GetGroupKey() keys.GroupKey {
	return m.Key.GroupKey
}

func (m *Membership) GetUserKey() keys.UserKey {
	return m.Key.UserKey
}

func (m *Membership) GetKey() keys.MembershipKey {
	return keys.NewMembershipKey(m.GetGroupKey(), m.GetUserKey())
}
