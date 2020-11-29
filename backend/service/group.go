package service

import (
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/group"
)

var _ group.Service = &GroupService{}

type GroupService struct {
	groupStore  group.Store
	amqpClient  amqp.Client
	chatService chat.Service
	authStore   auth.Store
}

func NewGroupService(groupStore group.Store, amqpClient amqp.Client, chatService chat.Service, authStore auth.Store) *GroupService {
	return &GroupService{
		groupStore:  groupStore,
		amqpClient:  amqpClient,
		chatService: chatService,
		authStore:   authStore,
	}
}
