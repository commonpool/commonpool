package group

import (
	"context"
	"github.com/commonpool/backend/model"
)

type Service interface {
	CreateGroup(ctx context.Context, request *CreateGroupRequest) (*CreateGroupResponse, error)
	GetGroup(ctx context.Context, request *GetGroupRequest) (*GetGroupResult, error)
	GetGroups(ctx context.Context, request *GetGroupsRequest) (*GetGroupsResult, error)
	GetGroupsByKeys(ctx context.Context, groupKeys []model.GroupKey) (*Groups, error)
	GetMembership(ctx context.Context, request *GetMembershipRequest) (*GetMembershipResponse, error)
	GetUserMemberships(ctx context.Context, request *GetMembershipsForUserRequest) (*GetMembershipsForUserResponse, error)
	GetGroupsMemberships(ctx context.Context, request *GetMembershipsForGroupRequest) (*GetMembershipsForGroupResponse, error)
	AcceptInvitation(ctx context.Context, request *AcceptInvitationRequest) (*AcceptInvitationResponse, error)
	DeclineInvitation(ctx context.Context, request *DeclineInvitationRequest) error
	LeaveGroup(ctx context.Context, request *LeaveGroupRequest) error
	SendGroupInvitation(ctx context.Context, request *InviteRequest) (*InviteResponse, error)
	RegisterUserAmqpSubscriptions(ctx context.Context) error
}

type AcceptInvitationRequest struct {
	MembershipKey model.MembershipKey
}

func NewAcceptInvitationRequest(membershipKey model.MembershipKey) *AcceptInvitationRequest {
	return &AcceptInvitationRequest{
		MembershipKey: membershipKey,
	}
}

type AcceptInvitationResponse struct {
	Membership *Membership
}

type DeclineInvitationRequest struct {
	MembershipKey model.MembershipKey
}

func NewDelineInvitationRequest(membershipKey model.MembershipKey) *DeclineInvitationRequest {
	return &DeclineInvitationRequest{
		MembershipKey: membershipKey,
	}
}

type LeaveGroupRequest struct {
	MembershipKey model.MembershipKey
}

func NewLeaveGroupRequest(membershipKey model.MembershipKey) *LeaveGroupRequest {
	return &LeaveGroupRequest{
		MembershipKey: membershipKey,
	}
}
