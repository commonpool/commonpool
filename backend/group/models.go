package group

import (
	"fmt"
	"github.com/commonpool/backend/model"
	"strconv"
	"time"
)

type Group struct {
	Key         model.GroupKey
	CreatedBy   model.UserKey
	CreatedAt   time.Time
	Name        string
	Description string
}

func (o *Group) GetKey() model.GroupKey {
	return o.Key
}

func (o *Group) GetCreatedByKey() model.UserKey {
	return o.CreatedBy
}

type Groups struct {
	Items []*Group
}

func NewGroups(groups []*Group) *Groups {
	return &Groups{
		Items: groups,
	}
}

type Names map[model.GroupKey]string

type PermissionType int

const (
	MemberPermission PermissionType = iota
	AdminPermission
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

type Memberships struct {
	Items []*Membership
}

func NewMemberships(items []*Membership) *Memberships {
	return &Memberships{Items: items}
}

func (m *Memberships) ContainsMembershipForGroup(groupKey model.GroupKey) bool {
	for _, item := range m.Items {
		if item.GetGroupKey().Equals(groupKey) {
			return true
		}
	}
	return false
}

type MembershipStatus int

const (
	ApprovedMembershipStatus MembershipStatus = iota
	PendingStatus
	PendingGroupMembershipStatus
	PendingUserMembershipStatus
)

func AnyMembershipStatus() *MembershipStatus {
	return nil
}

func ParseMembershipStatus(str string) (MembershipStatus, error) {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("cannot parse MembershipStatus: %s", err.Error())
	}
	if i < 0 || i > int(PendingUserMembershipStatus) {
		return 0, fmt.Errorf("cannot parse MembershipStatus")
	}
	return MembershipStatus(i), nil
}
