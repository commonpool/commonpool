package service

import (
	"context"
	group2 "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/service"
	"go.uber.org/zap"
)

func (g GroupService) GetUserMemberships(ctx context.Context, request *group2.GetMembershipsForUserRequest) (*group2.GetMembershipsForUserResponse, error) {

	ctx, l := service.GetCtx(ctx, "GroupService", "GetUserMemberships")

	l = l.With(zap.Object("user", request.UserKey))

	memberships, err := g.groupStore.GetMembershipsForUser(ctx, request.UserKey, request.MembershipStatus)
	if err != nil {
		l.Error("could not get memberships for user", zap.Error(err))
		return nil, err
	}

	return &group2.GetMembershipsForUserResponse{
		Memberships: memberships,
	}, nil
}
