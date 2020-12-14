package handler

import (
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/config"
	group2 "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/mq"
	resource2 "github.com/commonpool/backend/pkg/resource"
	trading2 "github.com/commonpool/backend/pkg/trading"
	"github.com/commonpool/backend/pkg/user"
)

type Handler struct {
	amqp           mq.Client
	resourceStore  resource2.Store
	authStore      user.Store
	authorization  auth.Authenticator
	chatStore      chat.Store
	tradingStore   trading2.Store
	groupService   group2.Service
	config         config.AppConfig
	chatService    chat.Service
	tradingService trading2.Service
}

func NewHandler(
	rs resource2.Store,
	as user.Store,
	cs chat.Store,
	ts trading2.Store,
	auth auth.Authenticator,
	amqp mq.Client,
	cfg config.AppConfig,
	chatService chat.Service,
	tradingService trading2.Service,
	groupService group2.Service,
) *Handler {
	return &Handler{
		resourceStore:  rs,
		authorization:  auth,
		authStore:      as,
		chatStore:      cs,
		tradingStore:   ts,
		config:         cfg,
		amqp:           amqp,
		chatService:    chatService,
		tradingService: tradingService,
		groupService:   groupService,
	}
}
