package handler

import (
	"github.com/commonpool/backend/pkg/auth"
	group "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/user"
	"github.com/labstack/echo/v4"
)

type GroupHandler struct {
	groupService group.Service
	userService  user.Service
	auth         auth.Authenticator
}

func NewHandler(groupService group.Service, userService user.Service, auth auth.Authenticator) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
		userService:  userService,
		auth:         auth,
	}
}

func (h *GroupHandler) Register(g *echo.Group) {

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
