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
		ID:          group.Key.String(),
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

func NewMembership(membership *group.Membership, groupNames group.Names, names auth.UserNames) Membership {
	return Membership{
		UserID:         membership.Key.UserKey.String(),
		GroupID:        membership.Key.GroupKey.String(),
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

type GetUserGroupsResponse struct {
	Groups []Group `json:"groups"`
}

type GetUserMembershipsResponse struct {
	Memberships []Membership `json:"memberships"`
}

func NewGetUserMembershipsResponse(memberships *group.Memberships, groupNames group.Names, userNames auth.UserNames) GetUserMembershipsResponse {
	responseMemberships := make([]Membership, len(memberships.Items))
	for i, membership := range memberships.Items {
		responseMemberships[i] = NewMembership(membership, groupNames, userNames)
	}
	return GetUserMembershipsResponse{
		Memberships: responseMemberships,
	}
}

type GetGroupMembershipsResponse struct {
	Memberships []Membership `json:"memberships"`
}

type GetUsersForGroupInvitePickerResponse struct {
	Users []UserInfoResponse `json:"users"`
	Take  int                `json:"take"`
	Skip  int                `json:"skip"`
}

type GetMembershipResponse struct {
	Membership Membership `json:"membership"`
}

func NewGetMembershipResponse(membership *group.Membership, groupNames group.Names, userNames auth.UserNames) *GetMembershipResponse {
	return &GetMembershipResponse{
		Membership: NewMembership(membership, groupNames, userNames),
	}
}

type CreateOrAcceptInvitationResponse struct {
	Membership Membership `json:"membership"`
}

func NewCreateOrAcceptInvitationResponse(membership *group.Membership, groupNames group.Names, userNames auth.UserNames) *CreateOrAcceptInvitationResponse {
	return &CreateOrAcceptInvitationResponse{
		Membership: NewMembership(membership, groupNames, userNames),
	}
}

type CreateOrAcceptInvitationRequest struct {
	UserID  string `json:"userId"`
	GroupID string `json:"groupId"`
}

type CancelOrDeclineInvitationResponse struct {
	Membership Membership `json:"membership"`
}

type CancelOrDeclineInvitationRequest struct {
	UserID  string `json:"userId"`
	GroupID string `json:"groupId"`
}
