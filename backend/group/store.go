package group

import (
	"context"
	"github.com/commonpool/backend/model"
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
	CreateGroupAndMembership(ctx context.Context, groupKey model.GroupKey, createdBy model.UserKey, name string, description string) (*Group, *Membership, error)
	GetGroup(ctx context.Context, groupKey model.GroupKey) (*Group, error)
	GetGroups(take int, skip int) (*Groups, int64, error)
	GetGroupsByKeys(ctx context.Context, groupKeys *model.GroupKeys) (*Groups, error)
	CreateMembership(ctx context.Context, membershipKey model.MembershipKey, isMember bool, isAdmin bool, isOwner bool, isDeactivated bool, groupConfirmed bool, userConfirmed bool) (*Membership, error)
	MarkInvitationAsAccepted(ctx context.Context, membershipKey model.MembershipKey, decisionFrom MembershipParty) error
	GetMembership(ctx context.Context, membershipKey model.MembershipKey) (*Membership, error)
	GetMembershipsForUser(ctx context.Context, userKey model.UserKey, membershipStatus *MembershipStatus) (*Memberships, error)
	GetMembershipsForGroup(ctx context.Context, groupKey model.GroupKey, membershipStatus *MembershipStatus) (*Memberships, error)
	DeleteMembership(ctx context.Context, membershipKey model.MembershipKey) error
}

type CreateGroupRequest struct {
	GroupKey    model.GroupKey
	Name        string
	Description string
}

type CreateGroupResponse struct {
	Group           *Group
	ChannelKey      model.ChannelKey
	Membership      *Membership
	SubscriptionKey model.ChannelSubscriptionKey
}

func NewCreateGroupRequest(key model.GroupKey, name string, description string) *CreateGroupRequest {
	return &CreateGroupRequest{
		GroupKey:    key,
		Name:        name,
		Description: description,
	}
}

type GetGroupRequest struct {
	Key model.GroupKey
}

type GetGroupResult struct {
	Group *Group
}

func NewGetGroupRequest(key model.GroupKey) *GetGroupRequest {
	return &GetGroupRequest{Key: key}
}

type MembershipPermissions struct {
	MembershipKey model.MembershipKey
	IsMember      bool
	IsAdmin       bool
}

type MembershipParty int

const (
	PartyGroup MembershipParty = iota
	PartyUser
)

type GetMembershipsForUserRequest struct {
	UserKey          model.UserKey
	MembershipStatus *MembershipStatus
}

func NewGetMembershipsForUserRequest(userKey model.UserKey, membershipStatus *MembershipStatus) *GetMembershipsForUserRequest {
	return &GetMembershipsForUserRequest{
		UserKey:          userKey,
		MembershipStatus: membershipStatus,
	}
}

type GetMembershipsForUserResponse struct {
	Memberships *Memberships
}

type GetMembershipsForGroupRequest struct {
	GroupKey         model.GroupKey
	MembershipStatus *MembershipStatus
}

func NewGetMembershipsForGroupRequest(groupKey model.GroupKey, status *MembershipStatus) *GetMembershipsForGroupRequest {
	return &GetMembershipsForGroupRequest{
		GroupKey:         groupKey,
		MembershipStatus: status,
	}
}

type GetMembershipsForGroupResponse struct {
	Memberships *Memberships
}

type GetMembershipRequest struct {
	MembershipKey model.MembershipKey
}

func NewGetMembershipRequest(membershipKey model.MembershipKey) *GetMembershipRequest {
	return &GetMembershipRequest{
		MembershipKey: membershipKey,
	}
}

type GetMembershipResponse struct {
	Membership *Membership
}
