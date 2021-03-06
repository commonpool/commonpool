package group

import (
	"context"
	"github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/keys"
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
	CreateGroupAndMembership(ctx context.Context, groupKey keys.GroupKey, createdBy keys.UserKey, name string, description string) (*Group, *domain.Membership, error)
	GetGroup(ctx context.Context, groupKey keys.GroupKey) (*Group, error)
	GetGroups(take int, skip int) (*Groups, int64, error)
	GetGroupsByKeys(ctx context.Context, groupKeys *keys.GroupKeys) (*Groups, error)
	CreateMembership(ctx context.Context, membershipKey keys.MembershipKey, isMember bool, isAdmin bool, isOwner bool, isDeactivated bool, groupConfirmed bool, userConfirmed bool) (*domain.Membership, error)
	MarkInvitationAsAccepted(ctx context.Context, membershipKey keys.MembershipKey, decisionFrom MembershipParty) error
	GetMembership(ctx context.Context, membershipKey keys.MembershipKey) (*domain.Membership, error)
	GetMembershipsForUser(ctx context.Context, userKey keys.UserKey, membershipStatus *domain.MembershipStatus) (*domain.Memberships, error)
	GetMembershipsForGroup(ctx context.Context, groupKey keys.GroupKey, membershipStatus *domain.MembershipStatus) (*domain.Memberships, error)
	DeleteMembership(ctx context.Context, membershipKey keys.MembershipKey) error
}

type CreateGroupRequest struct {
	GroupKey    keys.GroupKey
	Name        string
	Description string
}

type CreateGroupResponse struct {
	Group      *Group
	Membership *domain.Membership
}

func NewCreateGroupRequest(key keys.GroupKey, name string, description string) *CreateGroupRequest {
	return &CreateGroupRequest{
		GroupKey:    key,
		Name:        name,
		Description: description,
	}
}

type GetGroupRequest struct {
	Key keys.GroupKey
}

type GetGroupResult struct {
	Group *Group
}

func NewGetGroupRequest(key keys.GroupKey) *GetGroupRequest {
	return &GetGroupRequest{Key: key}
}

type MembershipPermissions struct {
	MembershipKey keys.MembershipKey
	IsMember      bool
	IsAdmin       bool
}

type MembershipParty int

const (
	PartyGroup MembershipParty = iota
	PartyUser
)

type GetMembershipsForUserRequest struct {
	UserKey          keys.UserKey
	MembershipStatus *domain.MembershipStatus
}

func NewGetMembershipsForUserRequest(userKey keys.UserKey, membershipStatus *domain.MembershipStatus) *GetMembershipsForUserRequest {
	return &GetMembershipsForUserRequest{
		UserKey:          userKey,
		MembershipStatus: membershipStatus,
	}
}

type GetMembershipsForUserResponse struct {
	Memberships *domain.Memberships
}

type GetMembershipsForGroupRequest struct {
	GroupKey         keys.GroupKey
	MembershipStatus *domain.MembershipStatus
}

func NewGetMembershipsForGroupRequest(groupKey keys.GroupKey, status *domain.MembershipStatus) *GetMembershipsForGroupRequest {
	return &GetMembershipsForGroupRequest{
		GroupKey:         groupKey,
		MembershipStatus: status,
	}
}

type GetMembershipsForGroupResponse struct {
	Memberships *domain.Memberships
}

type GetMembershipRequest struct {
	MembershipKey keys.MembershipKey
}

func NewGetMembershipRequest(membershipKey keys.MembershipKey) *GetMembershipRequest {
	return &GetMembershipRequest{
		MembershipKey: membershipKey,
	}
}

type GetMembershipResponse struct {
	Membership *domain.Membership
}
