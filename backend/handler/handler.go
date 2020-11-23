package handler

import (
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/config"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/store"
	"github.com/commonpool/backend/trading"
)

type Handler struct {
	resourceStore resource.Store
	authStore     auth.Store
	authorization auth.IAuth
	chatStore     chat.Store
	tradingStore  trading.Store
	groupStore    group.Store
	config        config.AppConfig
}

func NewHandler(
	rs resource.Store,
	as auth.Store,
	cs chat.Store,
	ts trading.Store,
	gs *store.GroupStore,
	auth auth.IAuth,
	cfg config.AppConfig,
) *Handler {
	return &Handler{
		resourceStore: rs,
		authorization: auth,
		authStore:     as,
		chatStore:     cs,
		tradingStore:  ts,
		groupStore:    gs,
		config:        cfg,
	}
}
