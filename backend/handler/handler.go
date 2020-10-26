package handler

import (
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/trading"
)

type Handler struct {
	resourceStore resource.Store
	authStore     auth.Store
	authorization auth.IAuth
	chatStore     chat.Store
	tradingStore  trading.Store
	//authStore    user.Store
	//articleStore article.Store
}

func NewHandler(rs resource.Store, as auth.Store, cs chat.Store, ts trading.Store, auth auth.IAuth) *Handler {
	return &Handler{
		resourceStore: rs,
		authorization: auth,
		authStore:     as,
		chatStore:     cs,
		tradingStore:  ts,
	}
}
