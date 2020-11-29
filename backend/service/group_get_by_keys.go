package service

import (
	"context"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"go.uber.org/zap"
)

func (g GroupService) GetGroupsByKeys(ctx context.Context, groupKeys []model.GroupKey) (*group.Groups, error) {

	ctx, l := GetCtx(ctx, "GroupService", "GetGroupsByKeys")

	groups, err := g.groupStore.GetGroupsByKeys(ctx, groupKeys)
	if err != nil {
		l.Error("could not get groups", zap.Error(err))
		return nil, err
	}

	return groups, nil
}
