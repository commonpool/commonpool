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

	c.Logger().Debug("SearchResources")

	// parsing request skip
	skip, err := utils.ParseSkip(c)
	if err != nil {
		c.Logger().Error(err, "SearchResource: failed to parse skip query param")
		return NewErrResponse(c, err)
	}

	// parsing request take
	take, err := utils.ParseTake(c, 0, 100)
	if err != nil {
		c.Logger().Error(err, "SearchResource: failed to parse take query param")
		return NewErrResponse(c, err)
	}

	// the search query
	searchQuery := strings.TrimSpace(c.QueryParam("query"))

	// the resource type
	resourceType, err := model.ParseResourceType(c.QueryParam("type"))
	if err != nil {
		c.Logger().Error(err, "SearchResource: failed to parse type query param")
		return NewErrResponse(c, err)
	}

	// created by
	createdBy := c.QueryParam("created_by")

	// visible in group
	var groupKey *model.GroupKey
	groupStr := c.QueryParam("group_id")
	if groupStr != "" {
		groupKey2, err := model.ParseGroupKey(groupStr)
		if err != nil {
			message := "SearchResource: could not parse group key"
			c.Logger().Error(err, message)
			return c.String(http.StatusInternalServerError, message)
		}
		groupKey = &groupKey2
	}

	// searching resources
	searchResourcesQuery := resource.NewSearchResourcesQuery(&searchQuery, resourceType, skip, take, createdBy, groupKey)
	searchResourcesResponse := h.resourceStore.Search(searchResourcesQuery)
	if searchResourcesResponse.Error != nil {
		c.Logger().Error(err, "SearchResource: failed to search resources from store")
		return NewErrResponse(c, searchResourcesResponse.Error)
	}

	// fetching groups with which resources are shared
	getGroupsResponse := h.groupStore.GetGroupsByKeys(group.NewGetGroupsByKeysQuery(searchResourcesResponse.Sharings.GetAllGroupKeys()))
	if getGroupsResponse.Error != nil {
		c.Logger().Error(err, "SearchResource: failed to get groups from store")
		return c.JSON(http.StatusInternalServerError, getGroupsResponse.Error.Error())
	}

	// building a map GroupKey -> Group for faster access
	groupMap := map[model.GroupKey]model.Group{}
	for _, g := range getGroupsResponse.Items {
		groupMap[g.GetKey()] = g
	}

	var createdByKeys []model.UserKey
	for _, item := range searchResourcesResponse.Resources.Items {
		createdByKeys = append(createdByKeys, item.GetUserKey())
	}

	createdByUsers, err := h.authStore.GetByKeys(createdByKeys)
	if err != nil {
		c.Logger().Error(err, "SearchResource: failed to get users by keys")
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// building response body
	resources := searchResourcesResponse.Resources.Items
	var resourcesResponse = make([]web.Resource, len(resources))
	for i, item := range resources {

		createdBy, err := createdByUsers.GetUser(model.NewUserKey(item.CreatedBy))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, getGroupsResponse.Error.Error())
		}

		// building list of groups that the resource is shared with
		var groups []model.Group
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
	c.Logger().Debug("GetResource: getting resource by id")

	resourceKey, err := model.ParseResourceKey(c.Param("id"))
	if err != nil {
		c.Logger().Error(err, "GetResource: could not parse resource key")
		return NewErrResponse(c, err)
	}

	getResourceByKeyResponse := h.resourceStore.GetByKey(resource.NewGetResourceByKeyQuery(*resourceKey))
	if getResourceByKeyResponse.Error != nil {
		c.Logger().Error(getResourceByKeyResponse.Error, "GetResource: could not get resource by key")
		return NewErrResponse(c, getResourceByKeyResponse.Error)
	}
	res := getResourceByKeyResponse.Resource

	createdBy := &model.User{}
	err = h.authStore.GetByKey(model.NewUserKey(res.CreatedBy), createdBy)
	if err != nil {
		c.Logger().Error(err, "GetResource: could not get user by key")
		return NewErrResponse(c, err)
	}

	// retrieving groups
	getGroupsResponse := h.groupStore.GetGroupsByKeys(group.NewGetGroupsByKeysQuery(getResourceByKeyResponse.Sharings.GetAllGroupKeys()))
	if getGroupsResponse.Error != nil {
		c.Logger().Error(getGroupsResponse.Error, "GetResource: could not get groups for resource")
		return NewErrResponse(c, getGroupsResponse.Error)
	}

	// return
	return c.JSON(http.StatusOK, web.GetResourceResponse{
		Resource: NewResourceResponse(res, createdBy.Username, createdBy.ID, getGroupsResponse.Items),
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
	c.Logger().Debug("CreateResource: creating resource")

	var err error

	// convert input body
	req := web.CreateResourceRequest{}
	if err = c.Bind(&req); err != nil {
		c.Logger().Error(err, "CreateResource: could not unmarshal request body")
		response := ErrCreateResourceBadRequest(err)
		return NewErrResponse(c, &response)
	}

	// validating body
	if err = c.Validate(req); err != nil {
		c.Logger().Error(err, "CreateResource: error validating request body")
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// sanitizing input
	sanitized := sanitizeCreateResource(req.Resource)

	// get logged in user
	loggedInUserSession := h.authorization.GetAuthUserSession(c)

	// creating new resource
	res := model.NewResource(
		model.NewResourceKey(uuid.NewV4()),
		sanitized.Type,
		loggedInUserSession.Subject,
		sanitized.Summary,
		sanitized.Description,
		sanitized.ValueInHoursFrom,
		sanitized.ValueInHoursTo,
	)

	// getting group keys that resource is shared with
	sharedWithGroupKeys, err, done := h.parseGroupKeys(c, req.Resource.SharedWith)
	if done {
		c.Logger().Error(err, "CreateResource: could not get group keys resource is shared with")
		return err
	}

	// making sure that user is active member of groups the resource is shared with
	err, done = h.ensureResourceIsSharedWithGroupsTheUserIsActiveMemberOf(c, loggedInUserSession.GetUserKey(), sharedWithGroupKeys)
	if done {
		c.Logger().Warn(err, "CreateResource: user tried to share resource with groups he's active member of")
		return err
	}

	// persist resource
	createResourceResponse := h.resourceStore.Create(resource.NewCreateResourceQuery(&res))
	if createResourceResponse.Error != nil {
		c.Logger().Error(createResourceResponse.Error, "CreateResource: could not persist resource")
		return NewErrResponse(c, createResourceResponse.Error)
	}

	// retrieving resource
	getResourceResponse := h.resourceStore.GetByKey(resource.NewGetResourceByKeyQuery(res.GetKey()))
	if getResourceResponse.Error != nil {
		c.Logger().Error(getResourceResponse.Error, "CreateResource: error while retrieving resource")
		return c.JSON(http.StatusBadRequest, getResourceResponse.Error)
	}

	// retrieving groups
	getGroupsResponse := h.groupStore.GetGroupsByKeys(group.NewGetGroupsByKeysQuery(getResourceResponse.Sharings.GetAllGroupKeys()))
	if getGroupsResponse.Error != nil {
		c.Logger().Error(getGroupsResponse.Error, "CreateResource: could not get groups")
		return NewErrResponse(c, getGroupsResponse.Error)
	}

	// send response
	return c.JSON(http.StatusCreated, web.CreateResourceResponse{
		Resource: NewResourceResponse(&res, loggedInUserSession.Username, loggedInUserSession.Subject, getGroupsResponse.Items),
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
		response := ErrValidation(err.Error())
		return NewErrResponse(c, &response)
	}

	// Gets the resource id
	resourceKey, err := model.ParseResourceKey(c.Param("id"))
	if err != nil {
		c.Logger().Error(err, "UpdateResource: error parsing resource key")
		response := ErrInvalidResourceKey(err.Error())
		return NewErrResponse(c, &response)
	}

	// Retrieves the resource
	getResourceByKeyResponse := h.resourceStore.GetByKey(resource.NewGetResourceByKeyQuery(*resourceKey))
	if getResourceByKeyResponse.Error != nil {
		c.Logger().Error(err, "UpdateResource: could not retrieve resource by key")
		return NewErrResponse(c, err)
	}

	resToUpdate := getResourceByKeyResponse.Resource

	// make sure user is owner of resource
	loggedInUser := h.authorization.GetAuthUserSession(c)
	if resToUpdate.GetUserKey() != loggedInUser.GetUserKey() {
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
		groupKey, err := model.ParseGroupKey(sharing.GroupID)
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
	getResourceResponse := h.resourceStore.GetByKey(resource.NewGetResourceByKeyQuery(resToUpdate.GetKey()))
	if getResourceResponse.Error != nil {
		c.Logger().Warn(err, "UpdateResource: error getting resource after update")
		return c.JSON(http.StatusBadRequest, getResourceResponse.Error)
	}

	// retrieving groups
	getGroupsResponse := h.groupStore.GetGroupsByKeys(group.NewGetGroupsByKeysQuery(getResourceResponse.Sharings.GetAllGroupKeys()))
	if getGroupsResponse.Error != nil {
		c.Logger().Warn(err, "UpdateResource: error getting groups")
		return c.JSON(http.StatusBadRequest, getGroupsResponse.Error)
	}

	return c.JSON(http.StatusOK, web.GetResourceResponse{
		Resource: NewResourceResponse(getResourceResponse.Resource, loggedInUser.Username, loggedInUser.Subject, getGroupsResponse.Items),
	})

}

func (h *Handler) ensureResourceIsSharedWithGroupsTheUserIsActiveMemberOf(c echo.Context, loggedInUserKey model.UserKey, sharedWithGroups []model.GroupKey) (error, bool) {
	var membershipStatus model.MembershipStatus = model.ApprovedMembershipStatus
	membershipsResponse := h.groupStore.GetMembershipsForUser(group.NewGetMembershipsForUserRequest(loggedInUserKey, &membershipStatus))

	// Checking if resource is shared with groups the user is part of
	for _, sharedWith := range sharedWithGroups {
		hasMembershipInGroup := membershipsResponse.Memberships.ContainsMembershipForGroup(sharedWith)
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
		groupKey, err := model.ParseGroupKey(groupKeyStr)
		if err != nil {
			return nil, c.String(http.StatusBadRequest, "invalid group key : "+groupKeyStr), true
		}
		sharedWithGroupKeys[i] = groupKey
	}
	return sharedWithGroupKeys, nil, false
}

func NewResourceResponse(res *model.Resource, creatorUsername string, creatorId string, sharedWithGroups []model.Group) web.Resource {

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
