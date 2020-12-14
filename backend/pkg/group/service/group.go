package service

import (
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/pkg/chat"
	group2 "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/user"
)

var _ group2.Service = &GroupService{}

type GroupService struct {
	groupStore  group2.Store
	amqpClient  amqp.Client
	chatService chat.Service
	authStore   user.Store
}

func NewGroupService(groupStore group2.Store, amqpClient amqp.Client, chatService chat.Service, authStore user.Store) *GroupService {
	return &GroupService{
		groupStore:  groupStore,
		amqpClient:  amqpClient,
		chatService: chatService,
		authStore:   authStore,
	}
}
