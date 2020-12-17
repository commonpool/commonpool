package service

import (
	"context"
	"github.com/commonpool/backend/pkg/group"
)

func (g GroupService) GetGroupMemberships(ctx context.Context, request *group.GetMembershipsForGroupRequest) (*group.GetMembershipsForGroupResponse, error) {

	memberships, err := g.groupStore.GetMembershipsForGroup(ctx, request.GroupKey, request.MembershipStatus)
	if err != nil {
		return nil, err
	}

	return &group.GetMembershipsForGroupResponse{
		Memberships: memberships,
	}, nil
}
