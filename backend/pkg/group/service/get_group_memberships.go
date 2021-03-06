package service

import (
	"context"
	"github.com/commonpool/backend/pkg/group"
)

func (g GroupService) GetGroupMemberships(ctx context.Context, request *group.GetMembershipsForGroupRequest) (*group.GetMembershipsForGroupResponse, error) {

	m, err := g.getGroupMemberships.Get(ctx, request.GroupKey, request.MembershipStatus)
	if err != nil {
		return nil, err
	}

	return &group.GetMembershipsForGroupResponse{
		Memberships: m,
	}, nil

}
