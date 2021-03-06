package service

import (
	"context"
	group2 "github.com/commonpool/backend/pkg/group"
)

func (g GroupService) GetMembership(ctx context.Context, request *group2.GetMembershipRequest) (*group2.GetMembershipResponse, error) {
	membership, err := g.groupStore.GetMembership(ctx, request.MembershipKey)
	if err != nil {
		return nil, err
	}
	return &group2.GetMembershipResponse{
		Membership: membership,
	}, nil

}
