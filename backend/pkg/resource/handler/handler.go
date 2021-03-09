package handler

import (
	"github.com/avast/retry-go"
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/auth/service"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/group/queries"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/domain"
	resourcequeries "github.com/commonpool/backend/pkg/resource/queries"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	"github.com/commonpool/backend/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type ResourceHandler struct {
	groupService                group.Service
	userService                 service.Service
	authorization               authenticator.Authenticator
	resourceRepo                domain.ResourceRepository
	getUserMemberships          *queries.GetUserMemberships
	getResource                 *resourcequeries.GetResource
	getResourceSharings         *resourcequeries.GetResourceSharings
	getResourcesSharings        *resourcequeries.GetResourcesSharings
	searchResources             *resourcequeries.SearchResources
	getResourceWithSharings     *resourcequeries.GetResourceWithSharings
	searchResourcesWithSharings *resourcequeries.SearchResourcesWithSharings
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
	searchResources *resourcequeries.SearchResources,
	getResourceWithSharings *resourcequeries.GetResourceWithSharings,
	searchResourcesWithSharings *resourcequeries.SearchResourcesWithSharings) *ResourceHandler {
	return &ResourceHandler{
		groupService:                groupService,
		userService:                 userService,
		authorization:               authenticator,
		resourceRepo:                resourceRepo,
		getUserMemberships:          getUserMemberships,
		getResource:                 getResource,
		getResourceSharings:         getResourceSharings,
		searchResources:             searchResources,
		getResourcesSharings:        getResourcesSharings,
		getResourceWithSharings:     getResourceWithSharings,
		searchResourcesWithSharings: searchResourcesWithSharings,
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

type GetResourceResponse struct {
	Resource *readmodel.ResourceWithSharingsReadModel `json:"resource"`
}

// GetResource
// @Summary Gets a single resource
// @Description Gets a resource by id
// @ID getResource
// @Tags resources
// @Accept json
// @Produce json
// @Param id path string true "Resource id" format(uuid)
// @Success 200 {object} web.GetResourceResponse
// @Failure 400 {object} utils.Error
// @Router /resources/:id [get]
func (h *ResourceHandler) GetResource(c echo.Context) error {
	ctx, _ := handler.GetEchoContext(c, "GetResource")
	resourceKey, err := keys.ParseResourceKey(c.Param("id"))
	if err != nil {
		return err
	}
	resource, err := h.getResourceWithSharings.Get(ctx, resourceKey)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, GetResourceResponse{
		Resource: resource,
	})
}

type SearchResourcesResponse struct {
	Take      int                                        `json:"take"`
	Skip      int                                        `json:"skip"`
	Resources []*readmodel.ResourceWithSharingsReadModel `json:"resources"`
}

// SearchResources godoc
// @Summary Searches resources
// @Description Search for resources
// @ID searchResources
// @Tags resources
// @Accept json
// @Produce json
// @Param query query string false "Search text"
// @Param type query string false "Resource type" Enums(0,1)
// @Param created_by query string false "Created by"
// @Param take query int false "Number of resources to take" minimum(0) maximum(100) default(10)
// @Param skip query int false "Number of resources to skip" minimum(0) default(0)
// @Success 200 {object} web.SearchResourcesResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 400 {object} utils.Error
// @Router /resources [get]
func (h *ResourceHandler) SearchResources(c echo.Context) error {

	ctx, l := handler.GetEchoContext(c, "SearchResources")
	l.Debug("searching resources")

	skip, err := utils.ParseSkip(c)
	if err != nil {
		return err
	}

	take, err := utils.ParseTake(c, 0, 100)
	if err != nil {
		return err
	}

	searchQuery := strings.TrimSpace(c.QueryParam("query"))

	var resourceType *domain.ResourceType
	resourceTypeStr := c.QueryParam("type")
	if resourceTypeStr != "" {
		resourceTypeValue, err := domain.ParseResourceType(resourceTypeStr)
		if err != nil {
			return err
		}
		resourceType = &resourceTypeValue
	}

	var callType *domain.CallType
	callTypeStr := c.QueryParam("sub_type")
	if callTypeStr != "" {
		callTypeValue, err := domain.ParseCallType(callTypeStr)
		if err != nil {
			return err
		}
		callType = &callTypeValue
	}

	var createdBy *string
	createdByStr := c.QueryParam("created_by")
	if createdByStr != "" {
		createdBy = &createdByStr
	}

	var groupKey *keys.GroupKey
	groupStr := c.QueryParam("group_id")
	if groupStr != "" {
		groupKey2, err := keys.ParseGroupKey(groupStr)
		if err != nil {
			return err
		}
		groupKey = &groupKey2
	}

	resourcesQuery := resourcequeries.NewSearchResourcesQuery(&searchQuery, resourceType, callType, skip, take, createdBy, groupKey)
	l.Debug("querying resources", zap.Object("query", resourcesQuery))

	resources, err := h.searchResourcesWithSharings.Get(ctx, resourcesQuery)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, SearchResourcesResponse{
		Resources: resources,
		Take:      take,
		Skip:      skip,
	})
}

type CreateResourceRequest struct {
	Resource CreateResourcePayload `json:"resource"`
}

type CreateResourcePayload struct {
	ResourceInfo domain.ResourceInfo   `json:"info"`
	SharedWith   InputResourceSharings `json:"sharings"`
}

type InputResourceSharing struct {
	GroupKey keys.GroupKey `json:"groupId" validate:"required,uuid"`
}

type InputResourceSharings []InputResourceSharing

func NewInputResourceSharings() InputResourceSharings {
	return InputResourceSharings{}
}

func (i InputResourceSharings) WithGroups(groupKeys ...keys.GroupKey) InputResourceSharings {
	var result = i
	for _, groupKey := range groupKeys {
		result = append(result, InputResourceSharing{
			GroupKey: groupKey,
		})
	}
	return result
}

func (p InputResourceSharings) GetGroupKeys() *keys.GroupKeys {
	sharedWithGroupKeys := make([]keys.GroupKey, len(p))
	for i := range p {
		sharedWithGroupKeys[i] = p[i].GroupKey
	}
	return keys.NewGroupKeys(sharedWithGroupKeys)
}

// CreateResource
// @Summary Creates a resource
// @Description Creates a resource
// @ID createResource
// @Tags resources
// @Accept json
// @Produce json
// @Param resource body web.CreateResourceRequest true "Resource to create"
// @Success 200 {object} web.CreateResourceResponse
// @Failure 400 {object} utils.Error
// @Router /resources [post]
func (h *ResourceHandler) CreateResource(c echo.Context) error {
	ctx, l := handler.GetEchoContext(c, "CreateResource")

	l.Debug("creating resource")

	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}

	req := CreateResourceRequest{}
	if err = c.Bind(&req); err != nil {
		return err
	}
	if err = c.Validate(req); err != nil {
		return err
	}

	resourceKey := keys.GenerateResourceKey()
	resource := domain.NewResource(resourceKey)

	err = resource.Register(
		loggedInUser.GetUserKey(),
		loggedInUser.GetUserKey().Target(),
		req.Resource.ResourceInfo,
		req.Resource.SharedWith.GetGroupKeys())
	if err != nil {
		l.Error("could not register resource", zap.Error(err))
		return err
	}

	if err := h.resourceRepo.Save(ctx, resource); err != nil {
		l.Error("could not save resource", zap.Error(err))
		return err
	}

	var rm *readmodel.ResourceWithSharingsReadModel
	err = retry.Do(func() error {
		rm, err = h.getResourceWithSharings.Get(ctx, resourceKey)
		if exceptions.Is(err, exceptions.ErrResourceNotFound) {
			return err
		}
		actualVersion := rm.Version
		expectedVersion := resource.GetVersion() - 1
		if actualVersion < expectedVersion {
			l.Debug("read model version not up to date", zap.Int("expected", expectedVersion), zap.Int("actual", actualVersion))
			return exceptions.ErrReadModelBackOff("resource", expectedVersion, actualVersion)
		}
		if err != nil {
			l.Warn("could not get resource read model", zap.Error(err))
			return err
		}
		return nil
	}, retry.Attempts(10), retry.MaxDelay(200*time.Millisecond))
	if err != nil {
		l.Error("failed to retrieve resource read model", zap.Error(err))
		return err
	}

	return c.JSON(http.StatusCreated, GetResourceResponse{
		Resource: rm,
	})
}

type UpdateResourceRequest struct {
	Resource UpdateResourcePayload `json:"resource"`
}

type UpdateResourcePayload struct {
	ResourceInfo domain.ResourceInfoUpdate `json:"info"`
	SharedWith   InputResourceSharings     `json:"sharedWith"`
}

// UpdateResource
// @Summary Updates a resource
// @Description Updates a resource
// @ID updateResource
// @Tags resources
// @Accept json
// @Produce json
// @Param id path string true "Resource id" format(uuid)
// @Param resource body web.UpdateResourceRequest true "Resource to create"
// @Success 200 {object} web.UpdateResourceResponse
// @Failure 400 {object} utils.Error
// @Router /resources [put]
func (h *ResourceHandler) UpdateResource(c echo.Context) error {

	ctx, l := handler.GetEchoContext(c, "UpdateResource")
	l = l.Named("ResourceHandler.UpdateResource")

	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}

	req := UpdateResourceRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	resourceKey, err := keys.ParseResourceKey(c.Param("id"))
	if err != nil {
		return err
	}

	resource, err := h.resourceRepo.Load(ctx, resourceKey)
	if err != nil {
		return err
	}
	if resource.GetVersion() == 0 {
		return exceptions.ErrResourceNotFound
	}

	err = resource.ChangeInfo(loggedInUser.GetUserKey(), req.Resource.ResourceInfo)
	if err != nil {
		return err
	}
	err = resource.ChangeSharings(loggedInUser.GetUserKey(), req.Resource.SharedWith.GetGroupKeys())

	if err := h.resourceRepo.Save(ctx, resource); err != nil {
		l.Error("could not save resource", zap.Error(err))
		return err
	}

	var rm *readmodel.ResourceWithSharingsReadModel
	err = retry.Do(func() error {
		rm, err = h.getResourceWithSharings.Get(ctx, resourceKey)
		if err != nil {
			return err
		}
		actualVersion := rm.Version
		expectedVersion := resource.GetVersion() - 1
		if rm.Version < expectedVersion {
			l.Debug("read model version not up to date", zap.Int("expected", expectedVersion), zap.Int("actual", actualVersion))
			return exceptions.ErrReadModelBackOff("resource", expectedVersion, actualVersion)
		}
		return nil
	}, retry.Attempts(10), retry.MaxDelay(200*time.Millisecond))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, GetResourceResponse{
		Resource: rm,
	})

}
