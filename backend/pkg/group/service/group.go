package service

import (
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/pkg/chat"
	group2 "github.com/commonpool/backend/pkg/group"
)

var _ group2.Service = &GroupService{}

type GroupService struct {
	groupStore  group2.Store
	amqpClient  amqp.Client
	chatService chat.Service
	authStore   auth.Store
}

func NewGroupService(groupStore group2.Store, amqpClient amqp.Client, chatService chat.Service, authStore auth.Store) *GroupService {
	return &GroupService{
		groupStore:  groupStore,
		amqpClient:  amqpClient,
		chatService: chatService,
		authStore:   authStore,
	}
}
