package service

import (
	"context"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/service"
	"go.uber.org/zap"
)

func (g GroupService) GetGroupsByKeys(ctx context.Context, groupKeys *model.GroupKeys) (*group.Groups, error) {

	ctx, l := service.GetCtx(ctx, "GroupService", "GetGroupsByKeys")

	if groupKeys == nil || len(groupKeys.Items) == 0 {
		return group.NewGroups([]*group.Group{}), nil
	}

	groups, err := g.groupStore.GetGroupsByKeys(ctx, groupKeys)
	if err != nil {
		l.Error("could not get groups", zap.Error(err))
		return nil, err
	}

	return groups, nil
}
