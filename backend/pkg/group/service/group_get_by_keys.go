package service

import (
	"context"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
)

func (g GroupService) GetGroupsByKeys(ctx context.Context, groupKeys *groupmodel.GroupKeys) (*groupmodel.Groups, error) {

	if groupKeys == nil || len(groupKeys.Items) == 0 {
		return groupmodel.NewGroups([]*groupmodel.Group{}), nil
	}

	groups, err := g.groupStore.GetGroupsByKeys(ctx, groupKeys)
	if err != nil {
		return nil, err
	}

	return groups, nil
}
