package service

import (
	"context"
	"github.com/commonpool/backend/pkg/group"
)

func (g GroupService) GetGroupsByKeys(ctx context.Context, groupKeys *group.GroupKeys) (*group.Groups, error) {

	if groupKeys == nil || len(groupKeys.Items) == 0 {
		return group.NewGroups([]*group.Group{}), nil
	}

	groups, err := g.groupStore.GetGroupsByKeys(ctx, groupKeys)
	if err != nil {
		return nil, err
	}

	return groups, nil
}
