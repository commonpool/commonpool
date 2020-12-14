package service

import (
	"context"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/service"
	"go.uber.org/zap"
)

func (g GroupService) GetUserMemberships(ctx context.Context, request *group.GetMembershipsForUserRequest) (*group.GetMembershipsForUserResponse, error) {

	ctx, l := service.GetCtx(ctx, "GroupService", "GetUserMemberships")

	l = l.With(zap.Object("user", request.UserKey))

	memberships, err := g.groupStore.GetMembershipsForUser(ctx, request.UserKey, request.MembershipStatus)
	if err != nil {
		l.Error("could not get memberships for user", zap.Error(err))
		return nil, err
	}

	return &group.GetMembershipsForUserResponse{
		Memberships: memberships,
	}, nil
}
