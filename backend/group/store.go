package group

import "github.com/commonpool/backend/model"

type GetGroupsByKeysQuery struct {
	GroupKeys []model.GroupKey
}

func NewGetGroupsByKeysQuery(groupKeys []model.GroupKey) *GetGroupsByKeysQuery {
	return &GetGroupsByKeysQuery{
		GroupKeys: groupKeys,
	}
}

type GetGroupsByKeysResponse struct {
	Items []model.Group
	Error error
}

func NewGetGroupsByKeysResponseSuccess(items []model.Group) *GetGroupsByKeysResponse {
	return &GetGroupsByKeysResponse{
		Items: items,
	}
}
func NewGetGroupsByKeysResponseError(err error) *GetGroupsByKeysResponse {
	return &GetGroupsByKeysResponse{
		Error: err,
	}
}

type Store interface {
	CreateGroup(request CreateGroupRequest) CreateGroupResponse
	GetGroup(request GetGroupRequest) GetGroupResult
	GetGroupsByKeys(request *GetGroupsByKeysQuery) *GetGroupsByKeysResponse
	GrantPermission(request GrantPermissionRequest) GrantPermissionResult
	RevokePermission(request RevokePermissionRequest) RevokePermissionResult
	Invite(request InviteRequest) InviteResponse
	Exclude(request ExcludeRequest) ExcludeResponse
	MarkInvitationAsAccepted(request MarkInvitationAsAcceptedRequest) MarkInvitationAsAcceptedResponse
	MarkInvitationAsDeclined(request MarkInvitationAsDeclinedRequest) MarkInvitationAsDeclinedResponse
	DeleteMembership(request DeleteMembershipRequest) DeleteMembershipResponse
	GetGroupPermissionsForUser(request GetMembershipPermissionsRequest) GetMembershipPermissionsResponse
	GetMembership(request GetMembershipRequest) GetMembershipResponse
	GetMembershipsForUser(request GetMembershipsForUserRequest) GetMembershipsForUserResponse
	GetMembershipsForGroup(request GetMembershipsForGroupRequest) GetMembershipsForGroupResponse
}

type CreateGroupRequest struct {
	GroupKey    model.GroupKey
	CreatedBy   model.UserKey
	Name        string
	Description string
}

type CreateGroupResponse struct {
	Error error
}

func NewCreateGroupRequest(key model.GroupKey, createdBy model.UserKey, name string, description string) CreateGroupRequest {
	return CreateGroupRequest{
		GroupKey:    key,
		CreatedBy:   createdBy,
		Name:        name,
		Description: description,
	}
}

type GetGroupRequest struct {
	Key model.GroupKey
}

type GetGroupResult struct {
	Error error
	Group model.Group
}

func NewGetGroupRequest(key model.GroupKey) GetGroupRequest {
	return GetGroupRequest{Key: key}
}

type GrantPermissionRequest struct {
	MembershipKey model.MembershipKey
	Permission    model.PermissionType
}

type GrantPermissionResult struct {
	Error error
}

func NewGrantPermissionRequest(membershipKey model.MembershipKey, permission model.PermissionType) GrantPermissionRequest {
	return GrantPermissionRequest{
		MembershipKey: membershipKey,
		Permission:    permission,
	}
}

type RevokePermissionRequest struct {
	MembershipKey model.MembershipKey
	Permission    model.PermissionType
}

type RevokePermissionResult struct {
	Error error
}

func NewRevokePermissionRequest(membershipKey model.MembershipKey, permission model.PermissionType) RevokePermissionRequest {
	return RevokePermissionRequest{
		MembershipKey: membershipKey,
		Permission:    permission,
	}
}

type InviteRequest struct {
	MembershipKey model.MembershipKey
	InvitedBy     MembershipParty
}

type InviteResponse struct {
	Error error
}

func NewInviteRequest(membershipKey model.MembershipKey) InviteRequest {
	return InviteRequest{
		MembershipKey: membershipKey,
	}
}

type ExcludeRequest struct {
	MembershipKey model.MembershipKey
}

type ExcludeResponse struct {
	Error error
}

func NewExcludeRequest(membershipKey model.MembershipKey) ExcludeRequest {
	return ExcludeRequest{
		MembershipKey: membershipKey,
	}
}

type GetMembershipPermissionsRequest struct {
	MembershipKey model.MembershipKey
}

type GetMembershipPermissionsResponse struct {
	Error                 error
	MembershipPermissions MembershipPermissions
}

func NewGetMembershipPermissionsRequest(membershipKey model.MembershipKey) GetMembershipPermissionsRequest {
	return GetMembershipPermissionsRequest{
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
	MembershipKey model.MembershipKey
	From          MembershipParty
}

func NewMarkInvitationAsAcceptedRequest(membershipKey model.MembershipKey, from MembershipParty) MarkInvitationAsAcceptedRequest {
	return MarkInvitationAsAcceptedRequest{MembershipKey: membershipKey, From: from}
}

type MarkInvitationAsAcceptedResponse struct {
	Error error
}

type MarkInvitationAsDeclinedRequest struct {
	MembershipKey model.MembershipKey
	From          MembershipParty
}

func NewMarkInvitationAsDeniedRequest(membershipKey model.MembershipKey, from MembershipParty) MarkInvitationAsDeclinedRequest {
	return MarkInvitationAsDeclinedRequest{MembershipKey: membershipKey, From: from}
}

type MarkInvitationAsDeclinedResponse struct {
	Error error
}

type GetMembershipsForUserRequest struct {
	UserKey          model.UserKey
	MembershipStatus *model.MembershipStatus
}

func NewGetMembershipsForUserRequest(userKey model.UserKey, membershipStatus *model.MembershipStatus) GetMembershipsForUserRequest {
	return GetMembershipsForUserRequest{
		UserKey:          userKey,
		MembershipStatus: membershipStatus,
	}
}

type GetMembershipsForUserResponse struct {
	Error       error
	Memberships model.Memberships
}

type GetMembershipsForGroupRequest struct {
	GroupKey         model.GroupKey
	MembershipStatus *model.MembershipStatus
}

func NewGetMembershipsForGroupRequest(groupKey model.GroupKey, status *model.MembershipStatus) GetMembershipsForGroupRequest {
	return GetMembershipsForGroupRequest{
		GroupKey:         groupKey,
		MembershipStatus: status,
	}
}

type GetMembershipsForGroupResponse struct {
	Error       error
	Memberships model.Memberships
}

type GetMembershipRequest struct {
	MembershipKey model.MembershipKey
}

func NewGetMembershipRequest(membershipKey model.MembershipKey) GetMembershipRequest {
	return GetMembershipRequest{
		MembershipKey: membershipKey,
	}
}

type GetMembershipResponse struct {
	Error      error
	Membership model.Membership
}

type DeleteMembershipRequest struct {
	MembershipKey model.MembershipKey
}

func NewDeleteMembershipRequest(membershipKey model.MembershipKey) DeleteMembershipRequest {
	return DeleteMembershipRequest{MembershipKey: membershipKey}
}

type DeleteMembershipResponse struct {
	Error error
}
