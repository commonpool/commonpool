package service

import (
	group2 "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/group/queries"
)

var _ group2.Service = &GroupService{}

type GroupService struct {
	groupRepo           domain.GroupRepository
	getByKeys           *queries.GetGroupByKeys
	getGroup            *queries.GetGroup
	getGroupMemberships *queries.GetGroupMemberships
}

func NewGroupService(
	groupRepo domain.GroupRepository,
	getGroup *queries.GetGroup,
	getByKeys *queries.GetGroupByKeys,
	getGroupMemberships *queries.GetGroupMemberships) *GroupService {
	return &GroupService{
		groupRepo:           groupRepo,
		getByKeys:           getByKeys,
		getGroup:            getGroup,
		getGroupMemberships: getGroupMemberships,
	}
}
