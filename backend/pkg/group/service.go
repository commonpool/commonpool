package group

import (
	"context"
	"github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/group/readmodels"
	"github.com/commonpool/backend/pkg/keys"
)

type Service interface {
	CreateGroup(ctx context.Context, request *CreateGroupRequest) (keys.GroupKey, error)
	GetGroup(ctx context.Context, key keys.GroupKey) (*readmodels.GroupReadModel, error)
	GetGroupsByKeys(ctx context.Context, groupKeys *keys.GroupKeys) ([]*readmodels.GroupReadModel, error)
	GetGroupMemberships(ctx context.Context, request *GetMembershipsForGroupRequest) (*GetMembershipsForGroupResponse, error)
	CreateOrAcceptInvitation(ctx context.Context, request *CreateOrAcceptInvitationRequest) error
	CancelOrDeclineInvitation(ctx context.Context, request *CancelOrDeclineInvitationRequest) error
}

type CreateOrAcceptInvitationRequest struct {
	MembershipKey keys.MembershipKey `json:"membership_key"`
}

func NewAcceptInvitationRequest(membershipKey keys.MembershipKey) *CreateOrAcceptInvitationRequest {
	return &CreateOrAcceptInvitationRequest{
		MembershipKey: membershipKey,
	}
}

type CreateOrAcceptInvitationResponse struct {
	Membership *domain.Membership
}

type CancelOrDeclineInvitationRequest struct {
	MembershipKey keys.MembershipKey `json:"membership_key"`
}

func NewDeclineInvitationRequest(membershipKey keys.MembershipKey) *CancelOrDeclineInvitationRequest {
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
