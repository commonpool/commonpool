package group

import (
	"context"
	chatmodel "github.com/commonpool/backend/pkg/chat/chatmodel"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
)

type GetGroupsRequest struct {
	Take int
	Skip int
}

type GetGroupsResult struct {
	Items      *groupmodel.Groups
	TotalCount int64
}

type Store interface {
	CreateGroupAndMembership(ctx context.Context, groupKey groupmodel.GroupKey, createdBy usermodel.UserKey, name string, description string) (*groupmodel.Group, *groupmodel.Membership, error)
	GetGroup(ctx context.Context, groupKey groupmodel.GroupKey) (*groupmodel.Group, error)
	GetGroups(take int, skip int) (*groupmodel.Groups, int64, error)
	GetGroupsByKeys(ctx context.Context, groupKeys *groupmodel.GroupKeys) (*groupmodel.Groups, error)
	CreateMembership(ctx context.Context, membershipKey groupmodel.MembershipKey, isMember bool, isAdmin bool, isOwner bool, isDeactivated bool, groupConfirmed bool, userConfirmed bool) (*groupmodel.Membership, error)
	MarkInvitationAsAccepted(ctx context.Context, membershipKey groupmodel.MembershipKey, decisionFrom MembershipParty) error
	GetMembership(ctx context.Context, membershipKey groupmodel.MembershipKey) (*groupmodel.Membership, error)
	GetMembershipsForUser(ctx context.Context, userKey usermodel.UserKey, membershipStatus *groupmodel.MembershipStatus) (*groupmodel.Memberships, error)
	GetMembershipsForGroup(ctx context.Context, groupKey groupmodel.GroupKey, membershipStatus *groupmodel.MembershipStatus) (*groupmodel.Memberships, error)
	DeleteMembership(ctx context.Context, membershipKey groupmodel.MembershipKey) error
}

type CreateGroupRequest struct {
	GroupKey    groupmodel.GroupKey
	Name        string
	Description string
}

type CreateGroupResponse struct {
	Group           *groupmodel.Group
	ChannelKey      chatmodel.ChannelKey
	Membership      *groupmodel.Membership
	SubscriptionKey chatmodel.ChannelSubscriptionKey
}

func NewCreateGroupRequest(key groupmodel.GroupKey, name string, description string) *CreateGroupRequest {
	return &CreateGroupRequest{
		GroupKey:    key,
		Name:        name,
		Description: description,
	}
}

type GetGroupRequest struct {
	Key groupmodel.GroupKey
}

type GetGroupResult struct {
	Group *groupmodel.Group
}

func NewGetGroupRequest(key groupmodel.GroupKey) *GetGroupRequest {
	return &GetGroupRequest{Key: key}
}

type MembershipPermissions struct {
	MembershipKey groupmodel.MembershipKey
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
	MembershipStatus *groupmodel.MembershipStatus
}

func NewGetMembershipsForUserRequest(userKey usermodel.UserKey, membershipStatus *groupmodel.MembershipStatus) *GetMembershipsForUserRequest {
	return &GetMembershipsForUserRequest{
		UserKey:          userKey,
		MembershipStatus: membershipStatus,
	}
}

type GetMembershipsForUserResponse struct {
	Memberships *groupmodel.Memberships
}

type GetMembershipsForGroupRequest struct {
	GroupKey         groupmodel.GroupKey
	MembershipStatus *groupmodel.MembershipStatus
}

func NewGetMembershipsForGroupRequest(groupKey groupmodel.GroupKey, status *groupmodel.MembershipStatus) *GetMembershipsForGroupRequest {
	return &GetMembershipsForGroupRequest{
		GroupKey:         groupKey,
		MembershipStatus: status,
	}
}

type GetMembershipsForGroupResponse struct {
	Memberships *groupmodel.Memberships
}

type GetMembershipRequest struct {
	MembershipKey groupmodel.MembershipKey
}

func NewGetMembershipRequest(membershipKey groupmodel.MembershipKey) *GetMembershipRequest {
	return &GetMembershipRequest{
		MembershipKey: membershipKey,
	}
}

type GetMembershipResponse struct {
	Membership *groupmodel.Membership
}
