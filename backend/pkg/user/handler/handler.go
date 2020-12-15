package handler

import (
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/user"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService   user.Service
	authenticator auth.Authenticator
}

func NewHandler(userService user.Service, authenticator auth.Authenticator) *UserHandler {
	return &UserHandler{
		userService:   userService,
		authenticator: authenticator,
	}
}

func (h *UserHandler) Register(g *echo.Group) {
	users := g.Group("/users", h.authenticator.Authenticate(true))
	users.GET("", h.SearchUsers)
	users.GET("/:id", h.GetUserInfo)
}
