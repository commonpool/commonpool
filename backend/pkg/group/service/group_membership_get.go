package service

import (
	"context"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/service"
	"go.uber.org/zap"
)

func (g GroupService) GetMembership(ctx context.Context, request *group.GetMembershipRequest) (*group.GetMembershipResponse, error) {

	ctx, l := service.GetCtx(ctx, "GroupService", "GetMembership")
	l = l.With(zap.Object("membership", request.MembershipKey))

	l.Debug("getting membership")

	membership, err := g.groupStore.GetMembership(ctx, request.MembershipKey)
	if err != nil {
		l.Error("could not get membership", zap.Error(err))
		return nil, err
	}

	return &group.GetMembershipResponse{
		Membership: membership,
	}, nil

}
