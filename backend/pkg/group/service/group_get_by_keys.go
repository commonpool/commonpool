package service

import (
	"context"
	"github.com/commonpool/backend/model"
	group2 "github.com/commonpool/backend/pkg/group"
)

func (g GroupService) GetGroupsByKeys(ctx context.Context, groupKeys *model.GroupKeys) (*group2.Groups, error) {

	if groupKeys == nil || len(groupKeys.Items) == 0 {
		return group2.NewGroups([]*group2.Group{}), nil
	}

	groups, err := g.groupStore.GetGroupsByKeys(ctx, groupKeys)
	if err != nil {
		return nil, err
	}

	return groups, nil
}
