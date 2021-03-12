package service

import (
	"context"
	"github.com/commonpool/backend/pkg/group/readmodels"
	"github.com/commonpool/backend/pkg/keys"
)

func (g *GroupService) GetGroupsByKeys(ctx context.Context, groupKeys *keys.GroupKeys) ([]*readmodels.GroupReadModel, error) {

	if groupKeys == nil || len(groupKeys.Items) == 0 {
		return []*readmodels.GroupReadModel{}, nil
	}

	return g.getByKeys.Get(groupKeys)

}
