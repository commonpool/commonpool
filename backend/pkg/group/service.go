package group

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
)

type Service interface {
	CreateGroup(ctx context.Context, request *CreateGroupRequest) (*CreateGroupResponse, error)
	GetGroup(ctx context.Context, request *GetGroupRequest) (*GetGroupResult, error)
	GetGroups(ctx context.Context, request *GetGroupsRequest) (*GetGroupsResult, error)
	GetGroupsByKeys(ctx context.Context, groupKeys *keys.GroupKeys) (*Groups, error)
	GetMembership(ctx context.Context, request *GetMembershipRequest) (*GetMembershipResponse, error)
	GetUserMemberships(ctx context.Context, request *GetMembershipsForUserRequest) (*GetMembershipsForUserResponse, error)
	GetGroupMemberships(ctx context.Context, request *GetMembershipsForGroupRequest) (*GetMembershipsForGroupResponse, error)
	CreateOrAcceptInvitation(ctx context.Context, request *CreateOrAcceptInvitationRequest) (*CreateOrAcceptInvitationResponse, error)
	CancelOrDeclineInvitation(ctx context.Context, request *CancelOrDeclineInvitationRequest) error
}

type CreateOrAcceptInvitationRequest struct {
	MembershipKey keys.MembershipKey
}

func NewAcceptInvitationRequest(membershipKey keys.MembershipKey) *CreateOrAcceptInvitationRequest {
	return &CreateOrAcceptInvitationRequest{
		MembershipKey: membershipKey,
	}
}

type CreateOrAcceptInvitationResponse struct {
	Membership *Membership
}

type CancelOrDeclineInvitationRequest struct {
	MembershipKey keys.MembershipKey
}

func NewDelineInvitationRequest(membershipKey keys.MembershipKey) *CancelOrDeclineInvitationRequest {
	return &CancelOrDeclineInvitationRequest{
		MembershipKey: membershipKey,
	}
}

type FindUsersForGroupInvitePickerQuery struct {
	Query    string
	Skip     int
	Take     int
	GroupKey keys.GroupKey
}
