package handler

import (
	"encoding/json"
	"github.com/avast/retry-go"
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/group/queries"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/domain"
	resourcequeries "github.com/commonpool/backend/pkg/resource/queries"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type ResourceHandler struct {
	groupService                     group.Service
	authorization                    authenticator.Authenticator
	resourceRepo                     domain.ResourceRepository
	getUserMemberships               *queries.GetUserMemberships
	getResource                      *resourcequeries.GetResource
	getResourceSharings              *resourcequeries.GetResourceSharings
	getResourcesSharings             *resourcequeries.GetResourcesSharings
	searchResources                  *resourcequeries.SearchResources
	getResourceWithSharings          *resourcequeries.GetResourceWithSharings
	getResourceWithSharingsAndValues *resourcequeries.GetResourceWithSharingsAndValues
	searchResourcesWithSharings      *resourcequeries.SearchResourcesWithSharings
	getUserResourceEvaluation        *resourcequeries.GetUserResourceEvaluation
}

func NewHandler(
	groupService group.Service,
	authenticator authenticator.Authenticator,
	resourceRepo domain.ResourceRepository,
	getUserMemberships *queries.GetUserMemberships,
	getResource *resourcequeries.GetResource,
	getResourceSharings *resourcequeries.GetResourceSharings,
	getResourcesSharings *resourcequeries.GetResourcesSharings,
	searchResources *resourcequeries.SearchResources,
	getResourceWithSharings *resourcequeries.GetResourceWithSharings,
	searchResourcesWithSharings *resourcequeries.SearchResourcesWithSharings,
	getResourceWithSharingsAndValues *resourcequeries.GetResourceWithSharingsAndValues,
	getUserResourceEvaluation *resourcequeries.GetUserResourceEvaluation,
) *ResourceHandler {
	return &ResourceHandler{
		groupService:                     groupService,
		authorization:                    authenticator,
		resourceRepo:                     resourceRepo,
		getUserMemberships:               getUserMemberships,
		getResource:                      getResource,
		getResourceSharings:              getResourceSharings,
		searchResources:                  searchResources,
		getResourcesSharings:             getResourcesSharings,
		getResourceWithSharings:          getResourceWithSharings,
		searchResourcesWithSharings:      searchResourcesWithSharings,
		getResourceWithSharingsAndValues: getResourceWithSharingsAndValues,
		getUserResourceEvaluation:        getUserResourceEvaluation,
	}
}

func (h *ResourceHandler) Register(e *echo.Group) {
	resources := e.Group("/resources", h.authorization.Authenticate(true))
	resources.GET("", h.SearchResources)
	resources.POST("", h.CreateResource)
	resources.GET("/:id", h.GetResource)
	resources.PUT("/:id", h.UpdateResource)
	resources.POST("/:id/inquire", h.InquireAboutResource)
	resources.PUT("/:id/evaluations", h.EvaluateResource)
	resources.GET("/:id/evaluations", h.GetMyEvaluation)
}

type GetResourceResponse struct {
	Resource *readmodel.ResourceWithSharingsAndValuesReadModel `json:"resource"`
}

func (g GetResourceResponse) GetResourceKey() keys.ResourceKey {
	return g.Resource.ResourceKey
}

func (g GetResourceResponse) AsUpdate() *UpdateResourceRequest {
	var sharedWith InputResourceSharings
	for _, sharing := range g.Resource.Sharings {
		sharedWith = append(sharedWith, InputResourceSharing{
			GroupKey: sharing.GroupKey,
		})
	}
	return &UpdateResourceRequest{
		Resource: UpdateResourcePayload{
			ResourceInfo: g.Resource.ResourceInfo.AsUpdate(),
			SharedWith:   sharedWith,
		},
	}
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
	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	resourceKey, err := keys.ParseResourceKey(c.Param("id"))
	if err != nil {
		return err
	}
	resource, err := h.getResourceWithSharingsAndValues.Get(ctx, resourceKey, loggedInUser.GetUserKey())
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

func (r SearchResourcesResponse) GetResourceKeys() []keys.ResourceKey {
	var result []keys.ResourceKey
	for _, resource := range r.Resources {
		result = append(result, resource.ResourceKey)
	}
	return result
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

	var query resourcequeries.SearchResourcesQuery
	if err := c.Bind(&query); err != nil {
		return exceptions.ErrBadRequest(err.Error())
	}
	if err := c.Validate(query); err != nil {
		return err
	}

	resources, err := h.searchResourcesWithSharings.Get(ctx, &query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, SearchResourcesResponse{
		Resources: resources,
		Take:      query.Take,
		Skip:      query.Skip,
	})
}

type CreateResourceRequest struct {
	Resource CreateResourcePayload `json:"resource"`
}

func NewCreateResourceRequest(resource CreateResourcePayload) *CreateResourceRequest {
	return &CreateResourceRequest{
		resource,
	}
}

type CreateResourcePayload struct {
	ResourceInfo domain.ResourceInfo     `json:"info"`
	SharedWith   InputResourceSharings   `json:"sharings"`
	Values       domain.ValueEstimations `json:"values"`
}

func NewCreateResourcePayload(resourceInfo domain.ResourceInfo, groupKeys ...keys.GroupKeyGetter) CreateResourcePayload {
	return CreateResourcePayload{}.WithResourceInfo(resourceInfo).SharedWithGroups(groupKeys...)
}

func (p CreateResourcePayload) WithResourceInfo(resourceInfo domain.ResourceInfo) CreateResourcePayload {
	return CreateResourcePayload{
		ResourceInfo: resourceInfo,
		SharedWith:   p.SharedWith,
	}
}

func (p CreateResourcePayload) SharedWithGroups(groupKeys ...keys.GroupKeyGetter) CreateResourcePayload {
	var sharings InputResourceSharings
	for _, groupKey := range groupKeys {
		sharings = append(sharings, InputResourceSharing{
			GroupKey: groupKey.GetGroupKey(),
		})
	}
	return CreateResourcePayload{
		ResourceInfo: p.ResourceInfo,
		SharedWith:   sharings,
	}
}

func (p CreateResourcePayload) AsRequest() *CreateResourceRequest {
	return &CreateResourceRequest{
		Resource: p,
	}
}

type InputResourceSharing struct {
	GroupKey keys.GroupKey `json:"groupId" validate:"required,uuid"`
}

type InputResourceSharings []InputResourceSharing

func (i *InputResourceSharings) UnmarshalJSON(data []byte) error {
	var res []InputResourceSharing
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	*i = res
	return nil
}

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
		req.Resource.SharedWith.GetGroupKeys(),
		req.Resource.Values)
	if err != nil {
		l.Error("could not register resource", zap.Error(err))
		return err
	}

	if err := h.resourceRepo.Save(ctx, resource); err != nil {
		l.Error("could not save resource", zap.Error(err))
		return err
	}

	var rm *readmodel.ResourceWithSharingsAndValuesReadModel
	err = retry.Do(func() error {
		rm, err = h.getResourceWithSharingsAndValues.Get(ctx, resourceKey, loggedInUser.GetUserKey())
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
	}, retry.Attempts(20), retry.MaxDelay(200*time.Millisecond))
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

func (u UpdateResourceRequest) WithShared(groups ...keys.GroupKeyGetter) *UpdateResourceRequest {
	var groupKeys []InputResourceSharing
	for _, getter := range groups {
		groupKeys = append(groupKeys, InputResourceSharing{
			GroupKey: getter.GetGroupKey(),
		})
	}
	return &UpdateResourceRequest{
		Resource: UpdateResourcePayload{
			ResourceInfo: u.Resource.ResourceInfo,
			SharedWith:   groupKeys,
		},
	}
}

type UpdateResourcePayload struct {
	ResourceInfo domain.ResourceInfoUpdate `json:"info"`
	SharedWith   InputResourceSharings     `json:"sharedWith"`
	Values       domain.ValueEstimations   `json:"values"`
}

func NewUpdateResourcePayload(resourceInfo domain.ResourceInfoUpdate, sharedWith ...keys.GroupKeyGetter) UpdateResourcePayload {
	var groupKeys []InputResourceSharing
	for _, getter := range sharedWith {
		groupKeys = append(groupKeys, InputResourceSharing{
			GroupKey: getter.GetGroupKey(),
		})
	}
	return UpdateResourcePayload{
		ResourceInfo: resourceInfo,
		SharedWith:   groupKeys,
	}
}

func (u UpdateResourcePayload) AsRequest() *UpdateResourceRequest {
	return &UpdateResourceRequest{
		Resource: u,
	}
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
	err = resource.EvaluateResource(loggedInUser.GetUserKey(), req.Resource.Values)

	if err := h.resourceRepo.Save(ctx, resource); err != nil {
		l.Error("could not save resource", zap.Error(err))
		return err
	}

	var rm *readmodel.ResourceWithSharingsAndValuesReadModel
	err = retry.Do(func() error {
		rm, err = h.getResourceWithSharingsAndValues.Get(ctx, resourceKey, loggedInUser.GetUserKey())
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

type GetUserResourceEvaluationsResponse struct {
	Evaluations []domain.DimensionValue `json:"values"`
}

func (h *ResourceHandler) GetMyEvaluation(c echo.Context) error {

	ctx, l := handler.GetEchoContext(c, "UpdateResource")
	l = l.Named("ResourceHandler.GetMyEvaluation")

	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}

	resourceKey, err := keys.ParseResourceKey(c.Param("id"))
	if err != nil {
		return err
	}

	evaluations, err := h.getUserResourceEvaluation.Get(ctx, resourceKey, loggedInUser.GetUserKey())
	if err != nil {
		return err
	}

	var response GetUserResourceEvaluationsResponse
	for _, evaluation := range evaluations {
		response.Evaluations = append(response.Evaluations, evaluation.DimensionValue)
	}

	return c.JSON(http.StatusOK, response)

}

type EvaluateResourceRequest struct {
	Values domain.ValueEstimations `json:"values"`
}

func (h *ResourceHandler) EvaluateResource(c echo.Context) error {

	ctx, l := handler.GetEchoContext(c, "UpdateResource")
	l = l.Named("ResourceHandler.GetMyEvaluation")

	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}

	resourceKey, err := keys.ParseResourceKey(c.Param("id"))
	if err != nil {
		return err
	}

	req := EvaluateResourceRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	resource, err := h.resourceRepo.Load(ctx, resourceKey)
	if err != nil {
		return err
	}

	if err := resource.EvaluateResource(loggedInUser.GetUserKey(), req.Values); err != nil {
		return err
	}

	if err := h.resourceRepo.Save(ctx, resource); err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)

}
