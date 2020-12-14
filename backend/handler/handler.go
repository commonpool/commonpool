package handler

import (
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/config"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/trading"
)

type Handler struct {
	amqp           amqp.Client
	resourceStore  resource.Store
	authStore      auth.Store
	authorization  auth.Authenticator
	chatStore      chat.Store
	tradingStore   trading.Store
	groupService   group.Service
	config         config.AppConfig
	chatService    chat.Service
	tradingService trading.Service
}

func NewHandler(
	rs resource.Store,
	as auth.Store,
	cs chat.Store,
	ts trading.Store,
	auth auth.Authenticator,
	amqp amqp.Client,
	cfg config.AppConfig,
	chatService chat.Service,
	tradingService trading.Service,
	groupService group.Service,
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
