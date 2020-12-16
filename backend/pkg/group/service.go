package group

import (
	"context"
)

type Service interface {
	CreateGroup(ctx context.Context, request *CreateGroupRequest) (*CreateGroupResponse, error)
	GetGroup(ctx context.Context, request *GetGroupRequest) (*GetGroupResult, error)
	GetGroups(ctx context.Context, request *GetGroupsRequest) (*GetGroupsResult, error)
	GetGroupsByKeys(ctx context.Context, groupKeys *GroupKeys) (*Groups, error)
	GetMembership(ctx context.Context, request *GetMembershipRequest) (*GetMembershipResponse, error)
	GetUserMemberships(ctx context.Context, request *GetMembershipsForUserRequest) (*GetMembershipsForUserResponse, error)
	GetGroupMemberships(ctx context.Context, request *GetMembershipsForGroupRequest) (*GetMembershipsForGroupResponse, error)
	CreateOrAcceptInvitation(ctx context.Context, request *CreateOrAcceptInvitationRequest) (*CreateOrAcceptInvitationResponse, error)
	CancelOrDeclineInvitation(ctx context.Context, request *CancelOrDeclineInvitationRequest) error
}

type CreateOrAcceptInvitationRequest struct {
	MembershipKey MembershipKey
}

func NewAcceptInvitationRequest(membershipKey MembershipKey) *CreateOrAcceptInvitationRequest {
	return &CreateOrAcceptInvitationRequest{
		MembershipKey: membershipKey,
	}
}

type CreateOrAcceptInvitationResponse struct {
	Membership *Membership
}

type CancelOrDeclineInvitationRequest struct {
	MembershipKey MembershipKey
}

func NewDelineInvitationRequest(membershipKey MembershipKey) *CancelOrDeclineInvitationRequest {
	return &CancelOrDeclineInvitationRequest{
		MembershipKey: membershipKey,
	}
}

type LeaveGroupRequest struct {
	MembershipKey MembershipKey
}
