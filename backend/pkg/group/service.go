package group

import (
	"context"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
)

type Service interface {
	CreateGroup(ctx context.Context, request *CreateGroupRequest) (*CreateGroupResponse, error)
	GetGroup(ctx context.Context, request *GetGroupRequest) (*GetGroupResult, error)
	GetGroups(ctx context.Context, request *GetGroupsRequest) (*GetGroupsResult, error)
	GetGroupsByKeys(ctx context.Context, groupKeys *groupmodel.GroupKeys) (*groupmodel.Groups, error)
	GetMembership(ctx context.Context, request *GetMembershipRequest) (*GetMembershipResponse, error)
	GetUserMemberships(ctx context.Context, request *GetMembershipsForUserRequest) (*GetMembershipsForUserResponse, error)
	GetGroupMemberships(ctx context.Context, request *GetMembershipsForGroupRequest) (*GetMembershipsForGroupResponse, error)
	CreateOrAcceptInvitation(ctx context.Context, request *CreateOrAcceptInvitationRequest) (*CreateOrAcceptInvitationResponse, error)
	CancelOrDeclineInvitation(ctx context.Context, request *CancelOrDeclineInvitationRequest) error
}

type CreateOrAcceptInvitationRequest struct {
	MembershipKey groupmodel.MembershipKey
}

func NewAcceptInvitationRequest(membershipKey groupmodel.MembershipKey) *CreateOrAcceptInvitationRequest {
	return &CreateOrAcceptInvitationRequest{
		MembershipKey: membershipKey,
	}
}

type CreateOrAcceptInvitationResponse struct {
	Membership *groupmodel.Membership
}

type CancelOrDeclineInvitationRequest struct {
	MembershipKey groupmodel.MembershipKey
}

func NewDelineInvitationRequest(membershipKey groupmodel.MembershipKey) *CancelOrDeclineInvitationRequest {
	return &CancelOrDeclineInvitationRequest{
		MembershipKey: membershipKey,
	}
}

type LeaveGroupRequest struct {
	MembershipKey groupmodel.MembershipKey
}
