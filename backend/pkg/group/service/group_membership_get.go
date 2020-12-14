package service

import (
	"context"
	group2 "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/service"
	"go.uber.org/zap"
)

func (g GroupService) GetMembership(ctx context.Context, request *group2.GetMembershipRequest) (*group2.GetMembershipResponse, error) {

	ctx, l := service.GetCtx(ctx, "GroupService", "GetMembership")
	l = l.With(zap.Object("membership", request.MembershipKey))

	l.Debug("getting membership")

	membership, err := g.groupStore.GetMembership(ctx, request.MembershipKey)
	if err != nil {
		l.Error("could not get membership", zap.Error(err))
		return nil, err
	}

	return &group2.GetMembershipResponse{
		Membership: membership,
	}, nil

}
