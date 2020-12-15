package web

import (
	"github.com/commonpool/backend/pkg/auth"
	model2 "github.com/commonpool/backend/pkg/group/model"
	"github.com/commonpool/backend/pkg/resource/model"
	"time"
)

type Group struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func NewGroup(group *model2.Group) *Group {
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

func NewMembership(membership *model2.Membership, groupNames model2.Names, names auth.UserNames) Membership {
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

func NewCreateGroupResponse(group *model2.Group) CreateGroupResponse {
	return CreateGroupResponse{
		Group: NewGroup(group),
	}
}

type GetGroupResponse struct {
	Group *Group `json:"group"`
}

func NewGetGroupResponse(group *model2.Group) GetGroupResponse {
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

func NewGetUserMembershipsResponse(memberships *model2.Memberships, groupNames model2.Names, userNames auth.UserNames) GetUserMembershipsResponse {
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

func NewGetMembershipResponse(membership *model2.Membership, groupNames model2.Names, userNames auth.UserNames) *GetMembershipResponse {
	return &GetMembershipResponse{
		Membership: NewMembership(membership, groupNames, userNames),
	}
}

type CreateOrAcceptInvitationResponse struct {
	Membership Membership `json:"membership"`
}

func NewCreateOrAcceptInvitationResponse(membership *model2.Membership, groupNames model2.Names, userNames auth.UserNames) *CreateOrAcceptInvitationResponse {
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

type OfferGroupOrUserPickerItem struct {
	Type    model.TargetType `json:"type"`
	UserID  *string          `json:"userId"`
	GroupID *string          `json:"groupId"`
	Name    string           `json:"name"`
}

type OfferGroupOrUserPickerResult struct {
	Items []OfferGroupOrUserPickerItem `json:"items"`
}
