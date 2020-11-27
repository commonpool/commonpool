package group

import (
	"context"
	"github.com/commonpool/backend/model"
)

type GetGroupsByKeysQuery struct {
	GroupKeys []model.GroupKey
}

func NewGetGroupsByKeysQuery(groupKeys []model.GroupKey) *GetGroupsByKeysQuery {
	return &GetGroupsByKeysQuery{
		GroupKeys: groupKeys,
	}
}

type GetGroupsByKeysResponse struct {
	Items []Group
}

type GetGroupsRequest struct {
	Take int
	Skip int
}

type GetGroupsResult struct {
	Items      []Group
	TotalCount int64
}

type Store interface {
	CreateGroup(ctx context.Context, groupKey model.GroupKey, createdBy model.UserKey, name string, description string) (*Group, error)
	GetGroup(ctx context.Context, groupKey model.GroupKey) (*Group, error)
	GetGroups(take int, skip int) ([]Group, int64, error)
	GetGroupsByKeys(ctx context.Context, groupKeys []model.GroupKey) (*Groups, error)
	GrantPermission(ctx context.Context, membershipKey model.MembershipKey, permission PermissionType) error
	RevokePermission(ctx context.Context, membershipKey model.MembershipKey, permission PermissionType) error
	CreateMembership(ctx context.Context, membershipKey model.MembershipKey, isMember bool, isAdmin bool, isOwner bool, isDeactivated bool, groupConfirmed bool, userConfirmed bool) (*Membership, error)
	Exclude(ctx context.Context, membershipKey model.MembershipKey) error
	MarkInvitationAsAccepted(ctx context.Context, membershipKey model.MembershipKey, decisionFrom MembershipParty) error
	MarkInvitationAsDeclined(ctx context.Context, membershipKey model.MembershipKey, decisionFrom MembershipParty) error
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

type GrantPermissionRequest struct {
	MembershipKey model.MembershipKey
	Permission    PermissionType
}

type GrantPermissionResult struct {
}

func NewGrantPermissionRequest(membershipKey model.MembershipKey, permission PermissionType) GrantPermissionRequest {
	return GrantPermissionRequest{
		MembershipKey: membershipKey,
		Permission:    permission,
	}
}

type RevokePermissionRequest struct {
	MembershipKey model.MembershipKey
	Permission    PermissionType
}

type RevokePermissionResult struct {
}

func NewRevokePermissionRequest(membershipKey model.MembershipKey, permission PermissionType) RevokePermissionRequest {
	return RevokePermissionRequest{
		MembershipKey: membershipKey,
		Permission:    permission,
	}
}

type InviteRequest struct {
	MembershipKey model.MembershipKey
}

type InviteResponse struct {
	Membership *Membership
}

func NewInviteRequest(membershipKey model.MembershipKey) *InviteRequest {
	return &InviteRequest{
		MembershipKey: membershipKey,
	}
}

type ExcludeRequest struct {
	MembershipKey model.MembershipKey
}

type ExcludeResponse struct {
}

func NewExcludeRequest(membershipKey model.MembershipKey) ExcludeRequest {
	return ExcludeRequest{
		MembershipKey: membershipKey,
	}
}

type MembershipPermissions struct {
	MembershipKey model.MembershipKey
	IsMember      bool
	IsAdmin       bool
}

type MembershipParty int

const (
	GroupParty MembershipParty = iota
	UserParty
)

type MarkInvitationAsAcceptedRequest struct {
	UserKey  model.UserKey
	GroupKey model.GroupKey
}

func NewMarkInvitationAsAcceptedRequest(groupKey model.GroupKey, userKey model.UserKey) *MarkInvitationAsAcceptedRequest {
	return &MarkInvitationAsAcceptedRequest{
		UserKey:  userKey,
		GroupKey: groupKey,
	}
}

type MarkInvitationAsAcceptedResponse struct {
}

type MarkInvitationAsDeclinedRequest struct {
	MembershipKey model.MembershipKey
	From          MembershipParty
}

func NewMarkInvitationAsDeniedRequest(membershipKey model.MembershipKey, from MembershipParty) MarkInvitationAsDeclinedRequest {
	return MarkInvitationAsDeclinedRequest{MembershipKey: membershipKey, From: from}
}

type MarkInvitationAsDeclinedResponse struct {
}

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

type DeleteMembershipRequest struct {
	MembershipKey model.MembershipKey
}

func NewDeleteMembershipRequest(membershipKey model.MembershipKey) DeleteMembershipRequest {
	return DeleteMembershipRequest{MembershipKey: membershipKey}
}

type DeleteMembershipResponse struct {
}
