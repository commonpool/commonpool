package service

import (
	"context"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/service"
	"go.uber.org/zap"
)

func (g GroupService) GetGroups(ctx context.Context, request *group.GetGroupsRequest) (*group.GetGroupsResult, error) {

	ctx, l := service.GetCtx(ctx, "GroupService", "GetGroups")

	l.Debug("getting groups")

	groups, totalCount, err := g.groupStore.GetGroups(request.Take, request.Skip)
	if err != nil {
		l.Error("could not get groups", zap.Error(err))
		return nil, err
	}

	return &group.GetGroupsResult{
		Items:      groups,
		TotalCount: totalCount,
	}, nil

}
