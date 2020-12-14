package service

import (
	"context"
	group2 "github.com/commonpool/backend/pkg/group"
)

func (g GroupService) GetGroup(ctx context.Context, request *group2.GetGroupRequest) (*group2.GetGroupResult, error) {

	grp, err := g.groupStore.GetGroup(ctx, request.Key)
	if err != nil {
		return nil, err
	}

	return &group2.GetGroupResult{
		Group: grp,
	}, nil

}
