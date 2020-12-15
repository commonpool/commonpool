package handler

import (
	"github.com/commonpool/backend/pkg/auth"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authorization auth.Authenticator
}

func NewHandler(authorization auth.Authenticator) *AuthHandler {
	return &AuthHandler{
		authorization: authorization,
	}
}

func (h *AuthHandler) Register(g *echo.Group) {
	grp := g.Group("/auth")
	grp.Any("/login", h.authorization.Login())
	grp.Any("/logout", h.authorization.Logout())
}
