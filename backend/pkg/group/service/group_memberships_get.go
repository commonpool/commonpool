package service

import (
	"context"
	group2 "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/service"
	"go.uber.org/zap"
)

func (g GroupService) GetGroupMemberships(ctx context.Context, request *group2.GetMembershipsForGroupRequest) (*group2.GetMembershipsForGroupResponse, error) {

	ctx, l := service.GetCtx(ctx, "GroupService", "GetGroupMemberships")

	l = l.With(zap.Object("group", request.GroupKey))

	memberships, err := g.groupStore.GetMembershipsForGroup(ctx, request.GroupKey, request.MembershipStatus)
	if err != nil {
		l.Error("could not get memberships for group", zap.Error(err))
		return nil, err
	}

	return &group2.GetMembershipsForGroupResponse{
		Memberships: memberships,
	}, nil
}
