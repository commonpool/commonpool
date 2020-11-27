package web

import (
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/group"
	"time"
)

type Group struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func NewGroup(group *group.Group) *Group {
	return &Group{
		ID:          group.ID.String(),
		CreatedAt:   group.CreatedAt,
		Name:        group.Name,
		Description: group.Description,
	}
}

type Membership struct {
	UserID         string    `json:"userId"`
	GroupID        string    `json:"groupId"`
	IsAdmin        bool      `json:"isAdmin"`
	IsMember       bool      `json:"isMember"`
	IsOwner        bool      `json:"isOwner"`
	GroupConfirmed bool      `json:"groupConfirmed"`
	UserConfirmed  bool      `json:"userConfirmed"`
	CreatedAt      time.Time `json:"createdAt"`
	IsDeactivated  bool      `json:"isDeactivated"`
	GroupName      string    `json:"groupName"`
	UserName       string    `json:"userName"`
}

func NewMembership(membership *group.Membership, groupNames group.GroupNames, names auth.UserNames) Membership {
	return Membership{
		UserID:         membership.UserID,
		GroupID:        membership.GroupID.String(),
		IsAdmin:        membership.IsAdmin,
		IsMember:       membership.IsMember,
		IsOwner:        membership.IsOwner,
		GroupConfirmed: membership.GroupConfirmed,
		UserConfirmed:  membership.UserConfirmed,
		CreatedAt:      membership.CreatedAt,
		IsDeactivated:  membership.IsDeactivated,
		GroupName:      groupNames[membership.GetGroupKey()],
		UserName:       names[membership.GetUserKey()],
	}
}

type CreateGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateGroupResponse struct {
	Group *Group `json:"group"`
}

func NewCreateGroupResponse(group *group.Group) CreateGroupResponse {
	return CreateGroupResponse{
		Group: NewGroup(group),
	}
}

type GetGroupResponse struct {
	Group *Group `json:"group"`
}

func NewGetGroupResponse(group *group.Group) GetGroupResponse {
	return GetGroupResponse{
		Group: NewGroup(group),
	}
}

type InviteUserRequest struct {
	UserID string `json:"userId"`
}

type InviteUserResponse struct {
	Membership Membership `json:"membership"`
}

func NewInviteUserResponse(membership *group.Membership, groupNames group.GroupNames, userNames auth.UserNames) InviteUserResponse {
	return InviteUserResponse{
		Membership: NewMembership(membership, groupNames, userNames),
	}
}

type ExcludeUserRequest struct {
	UserID string `json:"userId"`
}

type ExcludeUserResponse struct {
	Membership Membership `json:"membership"`
}

func NewExcludeUserResponse(membership group.Membership, groupNames group.GroupNames, userNames auth.UserNames) ExcludeUserResponse {
	return ExcludeUserResponse{
		Membership: NewMembership(&membership, groupNames, userNames),
	}
}

type GrantPermissionRequest struct {
	Permission group.PermissionType `json:"permission"`
	UserID     string               `json:"userId"`
	GroupID    string               `json:"groupId"`
}

type GrantPermissionResponse struct {
	Membership Membership `json:"membership"`
}

func NewGrantPermissionResponse(membership group.Membership, groupNames group.GroupNames, userNames auth.UserNames) GrantPermissionResponse {
	return GrantPermissionResponse{
		Membership: NewMembership(&membership, groupNames, userNames),
	}
}

type RevokePermissionRequest struct {
	Permission group.PermissionType `json:"permission"`
	UserID     string               `json:"userId"`
	GroupID    string               `json:"groupId"`
}

type RevokePermissionResponse struct {
	Membership Membership `json:"membership"`
}

func NewRevokePermissionResponse(membership group.Membership, groupNames group.GroupNames, userNames auth.UserNames) RevokePermissionResponse {
	return RevokePermissionResponse{
		Membership: NewMembership(&membership, groupNames, userNames),
	}
}

type GetUserGroupsResponse struct {
	Groups []Group `json:"groups"`
}

type GetUserMembershipsResponse struct {
	Memberships []Membership `json:"memberships"`
}

func NewGetUserMembershipsResponse(memberships *group.Memberships, groupNames group.GroupNames, userNames auth.UserNames) GetUserMembershipsResponse {
	responseMemberships := make([]Membership, len(memberships.Items))
	for i, membership := range memberships.Items {
		responseMemberships[i] = NewMembership(&membership, groupNames, userNames)
	}
	return GetUserMembershipsResponse{
		Memberships: responseMemberships,
	}
}

type GetGroupMembershipsResponse struct {
	Memberships []Membership `json:"memberships"`
}

func NewGetGroupMembershipsResponse(memberships []group.Membership, groupNames group.GroupNames, userNames auth.UserNames) GetGroupMembershipsResponse {
	responseMemberships := make([]Membership, len(memberships))
	for i, membership := range memberships {
		responseMemberships[i] = NewMembership(&membership, groupNames, userNames)
	}
	return GetGroupMembershipsResponse{
		Memberships: responseMemberships,
	}
}

type GetUsersForGroupInvitePickerResponse struct {
	Users []UserInfoResponse `json:"users"`
	Take  int                `json:"take"`
	Skip  int                `json:"skip"`
}

type GetMembershipResponse struct {
	Membership Membership `json:"membership"`
}

func NewGetMembershipResponse(membership *group.Membership, groupNames group.GroupNames, userNames auth.UserNames) GetMembershipResponse {
	return GetMembershipResponse{
		Membership: NewMembership(membership, groupNames, userNames),
	}
}

type AcceptInvitationResponse struct {
	Membership Membership `json:"membership"`
}

func NewAcceptInvitationResponse(membership *group.Membership, groupNames group.GroupNames, userNames auth.UserNames) *AcceptInvitationResponse {
	return &AcceptInvitationResponse{
		Membership: NewMembership(membership, groupNames, userNames),
	}
}

type DeclineInvitationResponse struct {
	Membership Membership `json:"membership"`
}

func NewDeclineInvitationResponse(membership group.Membership, groupNames group.GroupNames, userNames auth.UserNames) DeclineInvitationResponse {
	return DeclineInvitationResponse{
		Membership: NewMembership(&membership, groupNames, userNames),
	}
}

type LeaveGroupResponse struct {
	Membership Membership `json:"membership"`
}

func NewLeaveGroupResponse(membership group.Membership, groupNames group.GroupNames, userNames auth.UserNames) LeaveGroupResponse {
	return LeaveGroupResponse{
		Membership: NewMembership(&membership, groupNames, userNames),
	}
}
