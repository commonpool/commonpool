package handler

import (
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/auth/service"
	group "github.com/commonpool/backend/pkg/group"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	groupService group.Service
	userService  service.Service
	auth         authenticator.Authenticator
}

func NewHandler(groupService group.Service, userService service.Service, auth authenticator.Authenticator) *Handler {
	return &Handler{
		groupService: groupService,
		userService:  userService,
		auth:         auth,
	}
}

func (h *Handler) Register(g *echo.Group) {

	groups := g.Group("/groups", h.auth.Authenticate(true))
	groups.POST("", h.CreateGroup)
	groups.GET("/:id", h.GetGroup)
	groups.GET("/:id/memberships", h.GetGroupMemberships)
	groups.GET("/:id/memberships/:userId", h.GetMembership)
	groups.GET("/:id/invite-member-picker", h.GetUsersForGroupInvitePicker)

	memberships := g.Group("/memberships", h.auth.Authenticate(true))
	memberships.GET("", h.GetUserMemberships)
	memberships.POST("", h.CreateOrAcceptMembership)
	memberships.DELETE("", h.CancelOrDeclineInvitation)

}
