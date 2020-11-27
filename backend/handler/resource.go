package handler

import (
	"fmt"
	. "github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/utils"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

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
func (h *Handler) SearchResources(c echo.Context) error {

	ctx, l := GetEchoContext(c, "SearchResources")

	l.Debug("searching resources")

	l.Debug("parsing 'skip' query param")

	skip, err := utils.ParseSkip(c)
	if err != nil {
		l.Error("failed to parse skip query param", zap.Error(err))
		return NewErrResponse(c, err)
	}

	l.Debug("parsing 'take' query param")

	take, err := utils.ParseTake(c, 0, 100)
	if err != nil {
		l.Error("failed to parse take query param", zap.Error(err))
		return NewErrResponse(c, err)
	}

	l.Debug("parsing 'query' query param")

	searchQuery := strings.TrimSpace(c.QueryParam("query"))

	l.Debug("parsing 'type' query param")

	resourceType, err := resource.ParseResourceType(c.QueryParam("type"))
	if err != nil {
		l.Error("SearchResource: failed to parse type query param", zap.Error(err))
		return NewErrResponse(c, err)
	}

	l.Debug("parsing 'created_by' query param")

	createdBy := c.QueryParam("created_by")

	l.Debug("parsing 'group_id' query param")

	// visible in group
	var groupKey *model.GroupKey
	groupStr := c.QueryParam("group_id")
	if groupStr != "" {
		groupKey2, err := group.ParseGroupKey(groupStr)
		if err != nil {
			message := "SearchResource: could not parse group key"
			l.Error(message)
			return c.String(http.StatusInternalServerError, message)
		}
		groupKey = &groupKey2
	}

	l.Debug("searching resources")

	searchResourcesQuery := resource.NewSearchResourcesQuery(&searchQuery, resourceType, skip, take, createdBy, groupKey)
	searchResourcesResponse := h.resourceStore.Search(searchResourcesQuery)
	if searchResourcesResponse.Error != nil {
		c.Logger().Error(err, "SearchResource: failed to search resources from store")
		return NewErrResponse(c, searchResourcesResponse.Error)
	}

	l.Debug("fetching groups with which the resource is shared")

	getGroupsResponse, err := h.groupService.GetGroupsByKeys(ctx, searchResourcesResponse.Sharings.GetAllGroupKeys())
	if err != nil {
		l.Error("SearchResource: failed to get groups from store", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	l.Debug("building groupKey -> group map")

	groupMap := map[model.GroupKey]group.Group{}
	for _, g := range getGroupsResponse.Items {
		groupMap[g.GetKey()] = g
	}

	l.Debug("creating list of resource owners")

	var createdByKeys []model.UserKey
	for _, item := range searchResourcesResponse.Resources.Items {
		createdByKeys = append(createdByKeys, item.GetOwnerKey())
	}

	l.Debug("fetching resource owners")

	createdByUsers, err := h.authStore.GetByKeys(nil, createdByKeys)
	if err != nil {
		l.Error("failed to get users by keys", zap.Error(err))
		return c.String(http.StatusInternalServerError, err.Error())
	}

	l.Debug("building response body")

	// building response body
	resources := searchResourcesResponse.Resources.Items
	var resourcesResponse = make([]web.Resource, len(resources))
	for i, item := range resources {

		createdBy, err := createdByUsers.GetUser(model.NewUserKey(item.CreatedBy))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		// building list of groups that the resource is shared with
		var groups []group.Group
		sharings := searchResourcesResponse.Sharings.GetSharingsForResource(item.GetKey())
		for _, groupKey := range sharings.GetAllGroupKeys() {
			groups = append(groups, groupMap[groupKey])
		}

		// appending to result array
		resourcesResponse[i] = NewResourceResponse(&item, createdBy.Username, createdBy.ID, groups)
	}

	// return
	return c.JSON(http.StatusOK, web.SearchResourcesResponse{
		Resources:  resourcesResponse,
		Take:       take,
		Skip:       skip,
		TotalCount: searchResourcesResponse.TotalCount,
	})
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
func (h *Handler) GetResource(c echo.Context) error {

	ctx, l := GetEchoContext(c, "GetResource")

	l.Debug("parsing resource key")

	resourceKey, err := model.ParseResourceKey(c.Param("id"))
	if err != nil {
		l.Error("could not parse resource key", zap.Error(err))
		return NewErrResponse(c, err)
	}

	l.Debug("getting resource")

	getResourceByKeyResponse := h.resourceStore.GetByKey(ctx, resource.NewGetResourceByKeyQuery(*resourceKey))
	if getResourceByKeyResponse.Error != nil {
		l.Error("could not get resource by key", zap.Error(err))
		return NewErrResponse(c, getResourceByKeyResponse.Error)
	}
	res := getResourceByKeyResponse.Resource

	l.Debug("getting groups")

	groups, err := h.groupService.GetGroupsByKeys(ctx, getResourceByKeyResponse.Sharings.GetAllGroupKeys())
	if err != nil {
		l.Error("could not get groups", zap.Error(err))
		return err
	}

	ownerKey := res.GetOwnerKey()
	username, err := h.authStore.GetUsername(ownerKey)
	if err != nil {
		l.Error("could not get username", zap.Error(err))
		return err
	}

	// return
	return c.JSON(http.StatusOK, web.GetResourceResponse{
		Resource: NewResourceResponse(res, username, ownerKey.String(), groups.Items),
	})
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
func (h *Handler) CreateResource(c echo.Context) error {

	ctx, l := GetEchoContext(c, "CreateResource")

	var err error

	// convert input body
	req := web.CreateResourceRequest{}
	if err = c.Bind(&req); err != nil {
		l.Error("could not unmarshal request body", zap.Error(err))
		response := ErrCreateResourceBadRequest(err)
		return NewErrResponse(c, &response)
	}

	// validating body
	if err = c.Validate(req); err != nil {
		l.Error("CreateResource: error validating request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// sanitizing input
	sanitized := sanitizeCreateResource(req.Resource)

	// get logged in user
	loggedInUserSession := h.authorization.GetAuthUserSession(c)

	// getting group keys that resource is shared with
	sharedWithGroupKeys, err, done := h.parseGroupKeys(c, req.Resource.SharedWith)
	if done {
		l.Error("CreateResource: could not get group keys resource is shared with", zap.Error(err))
		return err
	}

	err, done = h.ensureResourceIsSharedWithGroupsTheUserIsActiveMemberOf(c, loggedInUserSession.GetUserKey(), sharedWithGroupKeys)
	if done {
		c.Logger().Warn(err, "CreateResource: user tried to share resource with groups he's active member of")
		return err
	}

	res := resource.NewResource(
		model.NewResourceKey(uuid.NewV4()),
		sanitized.Type,
		loggedInUserSession.Subject,
		sanitized.Summary,
		sanitized.Description,
		sanitized.ValueInHoursFrom,
		sanitized.ValueInHoursTo,
	)

	createResourceResponse := h.resourceStore.Create(resource.NewCreateResourceQuery(&res))
	if createResourceResponse.Error != nil {
		c.Logger().Error(createResourceResponse.Error, "CreateResource: could not persist resource")
		return NewErrResponse(c, createResourceResponse.Error)
	}

	getResourceResponse := h.resourceStore.GetByKey(ctx, resource.NewGetResourceByKeyQuery(res.GetKey()))
	if getResourceResponse.Error != nil {
		c.Logger().Error(getResourceResponse.Error, "CreateResource: error while retrieving resource")
		return c.JSON(http.StatusBadRequest, getResourceResponse.Error)
	}

	groups, err := h.groupService.GetGroupsByKeys(ctx, getResourceResponse.Sharings.GetAllGroupKeys())
	if err != nil {
		l.Error("CreateResource: could not get groups", zap.Error(err))
		return NewErrResponse(c, err)
	}

	// send response
	return c.JSON(http.StatusCreated, web.CreateResourceResponse{
		Resource: NewResourceResponse(&res, loggedInUserSession.Username, loggedInUserSession.Subject, groups.Items),
	})
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
func (h *Handler) UpdateResource(c echo.Context) error {

	ctx, l := GetEchoContext(c, "UpdateResource")

	c.Logger().Debug("UpdateResource: updating resource")
	req := web.UpdateResourceRequest{}

	// Binds the request payload to the web.UpdateResourceRequest instance
	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err, "UpdateResource: could not unmarshal request UpdateResourceRequest body")
		response := ErrUpdateResourceBadRequest(err)
		return NewErrResponse(c, &response)
	}

	// Validates the UpdateResourceRequest
	if err := c.Validate(req); err != nil {
		c.Logger().Warn(err, "UpdateResource: error validating updateResourceRequest")
		return NewErrResponse(c, ErrValidation(err.Error()))
	}

	// Gets the resource id
	resourceKey, err := model.ParseResourceKey(c.Param("id"))
	if err != nil {
		c.Logger().Error(err, "UpdateResource: error parsing resource key")
		response := ErrInvalidResourceKey(err.Error())
		return NewErrResponse(c, &response)
	}

	// Retrieves the resource
	getResourceByKeyResponse := h.resourceStore.GetByKey(ctx, resource.NewGetResourceByKeyQuery(*resourceKey))
	if getResourceByKeyResponse.Error != nil {
		c.Logger().Error(err, "UpdateResource: could not retrieve resource by key")
		return NewErrResponse(c, err)
	}

	resToUpdate := getResourceByKeyResponse.Resource

	// make sure user is owner of resource
	loggedInUser := h.authorization.GetAuthUserSession(c)
	if resToUpdate.GetOwnerKey() != loggedInUser.GetUserKey() {
		err := fmt.Errorf("cannot update a resource you do not own")
		c.Logger().Errorf("UpdateResource: %v", err)
		return c.String(http.StatusForbidden, err.Error())
	}

	// Parsing group keys that resource is shared with
	sharedWithGroupKeys, err, done := h.parseGroupKeys(c, req.Resource.SharedWith)
	if done {
		c.Logger().Error(err, "UpdateResource: could not parse group keys")
		return err
	}

	// make sure user is sharing resource with groups he's actively part of
	err, done = h.ensureResourceIsSharedWithGroupsTheUserIsActiveMemberOf(c, loggedInUser.GetUserKey(), sharedWithGroupKeys)
	if done {
		c.Logger().Warn(err, "UpdateResource: user tried to share resource with groups he's not actively part of")
		return err
	}

	// update resource
	req.Resource = sanitizeUpdateResource(req.Resource)
	resToUpdate.Summary = req.Resource.Summary
	resToUpdate.Description = req.Resource.Description
	resToUpdate.ValueInHoursFrom = req.Resource.ValueInHoursFrom
	resToUpdate.ValueInHoursTo = req.Resource.ValueInHoursTo

	// get shared with keys
	var groupKeys []model.GroupKey
	for _, sharing := range req.Resource.SharedWith {
		groupKey, err := group.ParseGroupKey(sharing.GroupID)
		if err != nil {
			message := "UpdateResource: could not parse groupKey"
			c.Logger().Error(err, message)
			return c.String(http.StatusInternalServerError, message)
		}
		groupKeys = append(groupKeys, groupKey)
	}

	// saving changes
	updateResourceQuery := resource.NewUpdateResourceQuery(resToUpdate, groupKeys)
	updateResourceResponse := h.resourceStore.Update(updateResourceQuery)
	if updateResourceResponse.Error != nil {
		c.Logger().Warn(updateResourceResponse.Error, "UpdateResource: error updating resource")
		return c.JSON(http.StatusBadRequest, "Could not update resource")
	}

	// retrieving resource
	getResourceResponse := h.resourceStore.GetByKey(ctx, resource.NewGetResourceByKeyQuery(resToUpdate.GetKey()))
	if getResourceResponse.Error != nil {
		c.Logger().Warn(err, "UpdateResource: error getting resource after update")
		return c.JSON(http.StatusBadRequest, getResourceResponse.Error)
	}

	// retrieving groups
	groups, err := h.groupService.GetGroupsByKeys(ctx, getResourceResponse.Sharings.GetAllGroupKeys())
	if err != nil {
		l.Error("could not get groups by keys", zap.Error(err))
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, web.GetResourceResponse{
		Resource: NewResourceResponse(getResourceResponse.Resource, loggedInUser.Username, loggedInUser.Subject, groups.Items),
	})

}

func (h *Handler) ensureResourceIsSharedWithGroupsTheUserIsActiveMemberOf(c echo.Context, loggedInUserKey model.UserKey, sharedWithGroups []model.GroupKey) (error, bool) {

	ctx, l := GetEchoContext(c, "ensureResourceIsSharedWithGroupsTheUserIsActiveMemberOf")

	var membershipStatus = group.ApprovedMembershipStatus

	userMemberships, err := h.groupService.GetUserMemberships(ctx, group.NewGetMembershipsForUserRequest(loggedInUserKey, &membershipStatus))
	if err != nil {
		l.Error("could not get user memberships", zap.Error(err))
		return err, true
	}

	// Checking if resource is shared with groups the user is part of
	for _, sharedWith := range sharedWithGroups {
		hasMembershipInGroup := userMemberships.Memberships.ContainsMembershipForGroup(sharedWith)
		if !hasMembershipInGroup {
			return c.String(http.StatusBadRequest, "cannot share resource with a group you are not part of"), true
		}
	}
	return nil, false
}

func sanitizeCreateResource(resource web.CreateResourcePayload) web.CreateResourcePayload {
	resource.Summary = strings.TrimSpace(resource.Summary)
	resource.Description = strings.TrimSpace(resource.Description)
	return resource
}

func sanitizeUpdateResource(resource web.UpdateResourcePayload) web.UpdateResourcePayload {
	resource.Summary = strings.TrimSpace(resource.Summary)
	resource.Description = strings.TrimSpace(resource.Description)
	return resource
}

func (h *Handler) parseGroupKeys(c echo.Context, sharedWith []web.InputResourceSharing) ([]model.GroupKey, error, bool) {
	sharedWithGroupKeys := make([]model.GroupKey, len(sharedWith))
	for i := range sharedWith {
		groupKeyStr := sharedWith[i].GroupID
		groupKey, err := group.ParseGroupKey(groupKeyStr)
		if err != nil {
			return nil, c.String(http.StatusBadRequest, "invalid group key : "+groupKeyStr), true
		}
		sharedWithGroupKeys[i] = groupKey
	}
	return sharedWithGroupKeys, nil, false
}

func NewResourceResponse(res *resource.Resource, creatorUsername string, creatorId string, sharedWithGroups []group.Group) web.Resource {

	//goland:noinspection GoPreferNilSlice
	var sharings = []web.OutputResourceSharing{}
	for _, withGroup := range sharedWithGroups {
		sharings = append(sharings, web.OutputResourceSharing{
			GroupID:   withGroup.ID.String(),
			GroupName: withGroup.Name,
		})
	}

	return web.Resource{
		Id:               res.ID.String(),
		Type:             res.Type,
		Description:      res.Description,
		Summary:          res.Summary,
		CreatedBy:        creatorUsername,
		CreatedById:      creatorId,
		CreatedAt:        res.CreatedAt,
		ValueInHoursFrom: res.ValueInHoursFrom,
		ValueInHoursTo:   res.ValueInHoursTo,
		SharedWith:       sharings,
	}
}

func NewErrResponse(c echo.Context, err error) error {
	res, ok := err.(*ErrorResponse)
	if !ok {
		statusCode := http.StatusInternalServerError
		return c.JSON(statusCode, NewError(err.Error(), "", statusCode))
	}
	return c.JSON(res.StatusCode, res)
}
