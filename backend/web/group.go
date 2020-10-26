package web

import (
	"github.com/commonpool/backend/model"
	"time"
)

type Group struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func NewGroup(group model.Group) Group {
	return Group{
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
	GroupConfirmed bool      `json:"groupConfirmed"`
	UserConfirmed  bool      `json:"userConfirmed"`
	CreatedAt      time.Time `json:"createdAt"`
	IsDeactivated  bool      `json:"isDeactivated"`
	GroupName      string    `json:"groupName"`
}

type GroupNames map[model.GroupKey]string

func NewMembership(membership model.Membership, groupNames GroupNames) Membership {
	return Membership{
		UserID:         membership.UserID,
		GroupID:        membership.GroupID.String(),
		IsAdmin:        membership.IsAdmin,
		IsMember:       membership.IsMember,
		GroupConfirmed: membership.GroupConfirmed,
		UserConfirmed:  membership.UserConfirmed,
		CreatedAt:      membership.CreatedAt,
		IsDeactivated:  membership.IsDeactivated,
		GroupName:      groupNames[membership.GetGroupKey()],
	}
}

type CreateGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateGroupResponse struct {
	Group Group `json:"group"`
}

func NewCreateGroupResponse(group model.Group) CreateGroupResponse {
	return CreateGroupResponse{
		Group: NewGroup(group),
	}
}

type GetGroupResponse struct {
	Group Group `json:"group"`
}

func NewGetGroupResponse(group model.Group) GetGroupResponse {
	return GetGroupResponse{
		Group: NewGroup(group),
	}
}

type InviteUserRequest struct {
	UserID  string `json:"userId"`
	GroupID string `json:"groupId"`
}

type InviteUserResponse struct {
	Membership Membership `json:"membership"`
}

func NewInviteUserResponse(membership model.Membership, groupNames GroupNames) InviteUserResponse {
	return InviteUserResponse{
		Membership: NewMembership(membership, groupNames),
	}
}

type ExcludeUserRequest struct {
	UserID  string `json:"userId"`
	GroupID string `json:"groupId"`
}

type ExcludeUserResponse struct {
	Membership Membership `json:"membership"`
}

func NewExcludeUserResponse(membership model.Membership, groupNames GroupNames) ExcludeUserResponse {
	return ExcludeUserResponse{
		Membership: NewMembership(membership, groupNames),
	}
}

type GrantPermissionRequest struct {
	Permission model.PermissionType `json:"permission"`
	UserID     string               `json:"userId"`
	GroupID    string               `json:"groupId"`
}

type GrantPermissionResponse struct {
	Membership Membership `json:"membership"`
}

func NewGrantPermissionResponse(membership model.Membership, groupNames GroupNames) GrantPermissionResponse {
	return GrantPermissionResponse{
		Membership: NewMembership(membership, groupNames),
	}
}

type RevokePermissionRequest struct {
	Permission model.PermissionType `json:"permission"`
	UserID     string               `json:"userId"`
	GroupID    string               `json:"groupId"`
}

type RevokePermissionResponse struct {
	Membership Membership `json:"membership"`
}

func NewRevokePermissionResponse(membership model.Membership, groupNames GroupNames) RevokePermissionResponse {
	return RevokePermissionResponse{
		Membership: NewMembership(membership, groupNames),
	}
}

type GetUserGroupsResponse struct {
	Groups []Group `json:"groups"`
}

type GetUserMembershipsResponse struct {
	Memberships []Membership `json:"memberships"`
}

func NewGetUserMembershipsResponse(memberships []model.Membership, groupNames GroupNames) GetUserMembershipsResponse {
	responseMemberships := make([]Membership, len(memberships))
	for i, membership := range memberships {
		responseMemberships[i] = NewMembership(membership, groupNames)
	}
	return GetUserMembershipsResponse{
		Memberships: responseMemberships,
	}
}
