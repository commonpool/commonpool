package service

import (
	"context"
	group2 "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/service"
	"go.uber.org/zap"
)

func (g GroupService) GetGroup(ctx context.Context, request *group2.GetGroupRequest) (*group2.GetGroupResult, error) {

	ctx, l := service.GetCtx(ctx, "GroupService", "GetGroup")
	l = l.With(zap.Object("group", request.Key))

	grp, err := g.groupStore.GetGroup(ctx, request.Key)
	if err != nil {
		l.Error("could not get group", zap.Error(err))
		return nil, err
	}

	return &group2.GetGroupResult{
		Group: grp,
	}, nil

}
