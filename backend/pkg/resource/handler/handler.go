package handler

import (
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/auth/service"
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/group/queries"
	"github.com/commonpool/backend/pkg/resource/domain"
	resourcequeries "github.com/commonpool/backend/pkg/resource/queries"
	"github.com/labstack/echo/v4"
)

type ResourceHandler struct {
	groupService         group.Service
	userService          service.Service
	authorization        authenticator.Authenticator
	resourceRepo         domain.ResourceRepository
	getUserMemberships   *queries.GetUserMemberships
	getResource          *resourcequeries.GetResource
	getResourceSharings  *resourcequeries.GetResourceSharings
	getResourcesSharings *resourcequeries.GetResourcesSharings
	searchResources      *resourcequeries.SearchResources
}

func NewHandler(
	groupService group.Service,
	userService service.Service,
	authenticator authenticator.Authenticator,
	resourceRepo domain.ResourceRepository,
	getUserMemberships *queries.GetUserMemberships,
	getResource *resourcequeries.GetResource,
	getResourceSharings *resourcequeries.GetResourceSharings,
	getResourcesSharings *resourcequeries.GetResourcesSharings,
	searchResources *resourcequeries.SearchResources) *ResourceHandler {
	return &ResourceHandler{
		groupService:         groupService,
		userService:          userService,
		authorization:        authenticator,
		resourceRepo:         resourceRepo,
		getUserMemberships:   getUserMemberships,
		getResource:          getResource,
		getResourceSharings:  getResourceSharings,
		searchResources:      searchResources,
		getResourcesSharings: getResourcesSharings,
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
