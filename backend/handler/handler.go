package handler

import (
	"github.com/commonpool/backend/resource"
)

type Handler struct {
	resourceStore resource.Store
	//userStore    user.Store
	//articleStore article.Store
}

func NewHandler(store resource.Store) *Handler {
	return &Handler{
		resourceStore: store,
		//userStore:    us,
		//articleStore: as,
	}
}
