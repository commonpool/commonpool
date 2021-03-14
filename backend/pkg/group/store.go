package group

import (
	"github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/group/readmodels"
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

type CreateGroupRequest struct {
	Name        string
	Description string
}

type CreateGroupResponse struct {
	Group      *Group
	Membership *domain.Membership
}

func NewCreateGroupRequest(name string, description string) *CreateGroupRequest {
	return &CreateGroupRequest{
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
	Memberships []*readmodels.MembershipReadModel
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
