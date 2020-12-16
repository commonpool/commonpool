package group

import (
	"context"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
)

type GetGroupsRequest struct {
	Take int
	Skip int
}

type GetGroupsResult struct {
	Items      *Groups
	TotalCount int64
}

type Store interface {
	CreateGroupAndMembership(ctx context.Context, groupKey GroupKey, createdBy usermodel.UserKey, name string, description string) (*Group, *Membership, error)
	GetGroup(ctx context.Context, groupKey GroupKey) (*Group, error)
	GetGroups(take int, skip int) (*Groups, int64, error)
	GetGroupsByKeys(ctx context.Context, groupKeys *GroupKeys) (*Groups, error)
	CreateMembership(ctx context.Context, membershipKey MembershipKey, isMember bool, isAdmin bool, isOwner bool, isDeactivated bool, groupConfirmed bool, userConfirmed bool) (*Membership, error)
	MarkInvitationAsAccepted(ctx context.Context, membershipKey MembershipKey, decisionFrom MembershipParty) error
	GetMembership(ctx context.Context, membershipKey MembershipKey) (*Membership, error)
	GetMembershipsForUser(ctx context.Context, userKey usermodel.UserKey, membershipStatus *MembershipStatus) (*Memberships, error)
	GetMembershipsForGroup(ctx context.Context, groupKey GroupKey, membershipStatus *MembershipStatus) (*Memberships, error)
	DeleteMembership(ctx context.Context, membershipKey MembershipKey) error
}

type CreateGroupRequest struct {
	GroupKey    GroupKey
	Name        string
	Description string
}

type CreateGroupResponse struct {
	Group      *Group
	Membership *Membership
}

func NewCreateGroupRequest(key GroupKey, name string, description string) *CreateGroupRequest {
	return &CreateGroupRequest{
		GroupKey:    key,
		Name:        name,
		Description: description,
	}
}

type GetGroupRequest struct {
	Key GroupKey
}

type GetGroupResult struct {
	Group *Group
}

func NewGetGroupRequest(key GroupKey) *GetGroupRequest {
	return &GetGroupRequest{Key: key}
}

type MembershipPermissions struct {
	MembershipKey MembershipKey
	IsMember      bool
	IsAdmin       bool
}

type MembershipParty int

const (
	PartyGroup MembershipParty = iota
	PartyUser
)

type GetMembershipsForUserRequest struct {
	UserKey          usermodel.UserKey
	MembershipStatus *MembershipStatus
}

func NewGetMembershipsForUserRequest(userKey usermodel.UserKey, membershipStatus *MembershipStatus) *GetMembershipsForUserRequest {
	return &GetMembershipsForUserRequest{
		UserKey:          userKey,
		MembershipStatus: membershipStatus,
	}
}

type GetMembershipsForUserResponse struct {
	Memberships *Memberships
}

type GetMembershipsForGroupRequest struct {
	GroupKey         GroupKey
	MembershipStatus *MembershipStatus
}

func NewGetMembershipsForGroupRequest(groupKey GroupKey, status *MembershipStatus) *GetMembershipsForGroupRequest {
	return &GetMembershipsForGroupRequest{
		GroupKey:         groupKey,
		MembershipStatus: status,
	}
}

type GetMembershipsForGroupResponse struct {
	Memberships *Memberships
}

type GetMembershipRequest struct {
	MembershipKey MembershipKey
}

func NewGetMembershipRequest(membershipKey MembershipKey) *GetMembershipRequest {
	return &GetMembershipRequest{
		MembershipKey: membershipKey,
	}
}

type GetMembershipResponse struct {
	Membership *Membership
}
