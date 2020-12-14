package service

import (
	"context"
	group2 "github.com/commonpool/backend/pkg/group"
)

func (g GroupService) GetGroupMemberships(ctx context.Context, request *group2.GetMembershipsForGroupRequest) (*group2.GetMembershipsForGroupResponse, error) {

	memberships, err := g.groupStore.GetMembershipsForGroup(ctx, request.GroupKey, request.MembershipStatus)
	if err != nil {
		return nil, err
	}

	return &group2.GetMembershipsForGroupResponse{
		Memberships: memberships,
	}, nil
}
