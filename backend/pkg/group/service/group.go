package service

import (
	"github.com/commonpool/backend/pkg/chat"
	group2 "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/group/queries"
	"github.com/commonpool/backend/pkg/mq"
	"github.com/commonpool/backend/pkg/user"
)

var _ group2.Service = &GroupService{}

type GroupService struct {
	groupStore  group2.Store
	amqpClient  mq.Client
	chatService chat.Service
	authStore   user.Store
	groupRepo   domain.GroupRepository
	getByKeys   *queries.GetGroupByKeys
	getGroup    *queries.GetGroup
}

func NewGroupService(
	groupStore group2.Store,
	amqpClient mq.Client,
	chatService chat.Service,
	authStore user.Store,
	groupRepo domain.GroupRepository,
	getGroup *queries.GetGroup,
	getByKeys *queries.GetGroupByKeys) *GroupService {
	return &GroupService{
		groupStore:  groupStore,
		amqpClient:  amqpClient,
		chatService: chatService,
		authStore:   authStore,
		groupRepo:   groupRepo,
		getByKeys:   getByKeys,
		getGroup:    getGroup,
	}
}
