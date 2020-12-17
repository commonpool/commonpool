package service

import (
	"context"
	group2 "github.com/commonpool/backend/pkg/group"
)

func (g GroupService) GetGroups(ctx context.Context, request *group2.GetGroupsRequest) (*group2.GetGroupsResult, error) {

	groups, totalCount, err := g.groupStore.GetGroups(request.Take, request.Skip)
	if err != nil {
		return nil, err
	}

	return &group2.GetGroupsResult{
		Items:      groups,
		TotalCount: totalCount,
	}, nil

}
