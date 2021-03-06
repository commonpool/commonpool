package handler

import (
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/auth/service"
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/group/queries"
	"github.com/commonpool/backend/pkg/resource"
	"github.com/labstack/echo/v4"
)

type ResourceHandler struct {
	resourceService    resource.Service
	groupService       group.Service
	userService        service.Service
	authorization      authenticator.Authenticator
	getUserMemberships *queries.GetUserMemberships
}

func NewHandler(
	resourceService resource.Service,
	groupService group.Service,
	userService service.Service,
	authenticator authenticator.Authenticator,
	getUserMemberships *queries.GetUserMemberships) *ResourceHandler {
	return &ResourceHandler{
		resourceService:    resourceService,
		groupService:       groupService,
		userService:        userService,
		authorization:      authenticator,
		getUserMemberships: getUserMemberships,
	}
}

func (h *ResourceHandler) Register(e *echo.Group) {
	resources := e.Group("/resources", h.authorization.Authenticate(true))
	resources.GET("", h.SearchResources)
	resources.POST("", h.CreateResource)
	resources.GET("/:id", h.GetResource)
	resources.PUT("/:id", h.UpdateResource)
	resources.POST("/:id/inquire", h.InquireAboutResource)
}
