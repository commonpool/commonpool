package handler

import (
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/resource"
)

type Handler struct {
	resourceStore resource.Store
	userStore     auth.Store
	authorization auth.IAuth
	//userStore    user.Store
	//articleStore article.Store
}

func NewHandler(store resource.Store, authStore auth.Store, authorization auth.IAuth) *Handler {
	return &Handler{
		resourceStore: store,
		authorization: authorization,
		userStore:     authStore,
		//userStore:    us,
		//articleStore: as,
	}
}
