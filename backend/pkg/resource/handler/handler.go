package handler

import (
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/resource"
	"github.com/commonpool/backend/pkg/user"
	"github.com/labstack/echo/v4"
)

type ResourceHandler struct {
	resourceService resource.Service
	groupService    group.Service
	userService     user.Service
	authorization   auth.Authenticator
}

func NewHandler(resourceService resource.Service, groupService group.Service, userService user.Service, authenticator auth.Authenticator) *ResourceHandler {
	return &ResourceHandler{
		resourceService: resourceService,
		groupService:    groupService,
		userService:     userService,
		authorization:   authenticator,
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
