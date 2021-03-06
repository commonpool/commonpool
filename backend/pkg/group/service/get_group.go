package service

import (
	"context"
	"github.com/commonpool/backend/pkg/group/readmodels"
	"github.com/commonpool/backend/pkg/keys"
)

func (g GroupService) GetGroup(ctx context.Context, groupKey keys.GroupKey) (*readmodels.GroupReadModel, error) {
	group, err := g.getGroup.Get(ctx, groupKey)
	if err != nil {
		return nil, err
	}
	return group, nil
}
