package service

import (
	"context"
	"github.com/commonpool/backend/group"
	"go.uber.org/zap"
)

func (g GroupService) GetGroupsMemberships(ctx context.Context, request *group.GetMembershipsForGroupRequest) (*group.GetMembershipsForGroupResponse, error) {

	ctx, l := GetCtx(ctx, "GroupService", "GetGroupsMemberships")

	l = l.With(zap.Object("group", request.GroupKey))

	memberships, err := g.groupStore.GetMembershipsForGroup(ctx, request.GroupKey, request.MembershipStatus)
	if err != nil {
		l.Error("could not get memberships for group", zap.Error(err))
		return nil, err
	}

	return &group.GetMembershipsForGroupResponse{
		Memberships: memberships,
	}, nil
}
