package group

import (
	"fmt"
	"github.com/commonpool/backend/model"
	"github.com/satori/go.uuid"
	"strconv"
	"time"
)

type Group struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	CreatedBy   string
	CreatedAt   time.Time
	Name        string
	Description string
}

func (o *Group) GetKey() model.GroupKey {
	return model.NewGroupKey(o.ID)
}

func (o *Group) GetCreatedByKey() model.UserKey {
	return model.NewUserKey(o.CreatedBy)
}

type Groups struct {
	Items []Group
}

func NewGroups(groups []Group) *Groups {
	return &Groups{
		Items: groups,
	}
}

type GroupNames map[model.GroupKey]string

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
	membershipKey model.MembershipKey,
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

func (m *Membership) GetGroupKey() model.GroupKey {
	return model.NewGroupKey(m.GroupID)
}

func (m *Membership) GetUserKey() model.UserKey {
	return model.NewUserKey(m.UserID)
}

func (m *Membership) GetKey() model.MembershipKey {
	return model.NewMembershipKey(m.GetGroupKey(), m.GetUserKey())
}

type Memberships struct {
	Items []Membership
}

func NewMemberships(items []Membership) *Memberships {
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
