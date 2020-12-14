package service

import (
	"context"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/service"
	"go.uber.org/zap"
)

func (g GroupService) GetGroup(ctx context.Context, request *group.GetGroupRequest) (*group.GetGroupResult, error) {

	ctx, l := service.GetCtx(ctx, "GroupService", "GetGroup")
	l = l.With(zap.Object("group", request.Key))

	grp, err := g.groupStore.GetGroup(ctx, request.Key)
	if err != nil {
		l.Error("could not get group", zap.Error(err))
		return nil, err
	}

	return &group.GetGroupResult{
		Group: grp,
	}, nil

}
