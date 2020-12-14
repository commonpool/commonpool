package handler

import (
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/config"
	"github.com/commonpool/backend/pkg/chat"
	group2 "github.com/commonpool/backend/pkg/group"
	resource2 "github.com/commonpool/backend/pkg/resource"
	trading2 "github.com/commonpool/backend/pkg/trading"
)

type Handler struct {
	amqp           amqp.Client
	resourceStore  resource2.Store
	authStore      auth.Store
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
	as auth.Store,
	cs chat.Store,
	ts trading2.Store,
	auth auth.Authenticator,
	amqp amqp.Client,
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
