package handler

import (
	"fmt"
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/exceptions"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	"github.com/commonpool/backend/pkg/handler"
	resource "github.com/commonpool/backend/pkg/resource"
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

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

	c.Logger().Debug("UpdateResource: updating resource")
	req := web.UpdateResourceRequest{}

	// Binds the request payload to the web.UpdateResourceRequest instance
	if err := c.Bind(&req); err != nil {
		return err
	}

	// Validates the UpdateResourceRequest
	if err := c.Validate(req); err != nil {
		return err
	}

	// Gets the resource id
	resourceKey, err := resourcemodel.ParseResourceKey(c.Param("id"))
	if err != nil {
		return err
	}

	// Retrieves the resource
	getResourceByKeyResponse, err := h.resourceService.GetByKey(ctx, resource.NewGetResourceByKeyQuery(resourceKey))
	if err != nil {
		return err
	}

	resToUpdate := getResourceByKeyResponse.Resource

	// make sure user is owner of resource

	loggedInUser, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return exceptions.ErrUnauthorized
	}
	if resToUpdate.GetOwnerKey() != loggedInUser.GetUserKey() {
		err := fmt.Errorf("cannot update a resource you do not own")
		return err
	}

	// Parsing group keys that resource is shared with
	sharedWithGroupKeys, err, done := h.parseGroupKeys(c, req.Resource.SharedWith)
	if done {
		return err
	}

	// make sure user is sharing resource with groups he's actively part of
	err, done = h.ensureResourceIsSharedWithGroupsTheUserIsActiveMemberOf(c, loggedInUser.GetUserKey(), sharedWithGroupKeys)
	if done {
		return err
	}

	// update resource
	req.Resource.Summary = strings.TrimSpace(req.Resource.Summary)
	req.Resource.Description = strings.TrimSpace(req.Resource.Description)
	resToUpdate.Summary = req.Resource.Summary
	resToUpdate.Description = req.Resource.Description
	resToUpdate.ValueInHoursFrom = req.Resource.ValueInHoursFrom
	resToUpdate.ValueInHoursTo = req.Resource.ValueInHoursTo

	// get shared with keys
	var groupKeys []groupmodel.GroupKey
	for _, sharing := range req.Resource.SharedWith {
		groupKey, err := groupmodel.ParseGroupKey(sharing.GroupID)
		if err != nil {
			message := "UpdateResource: could not parse groupKey"
			c.Logger().Error(err, message)
			return c.String(http.StatusInternalServerError, message)
		}
		groupKeys = append(groupKeys, groupKey)
	}

	if err := h.resourceService.Update(ctx, resource.NewUpdateResourceQuery(resToUpdate, groupmodel.NewGroupKeys(groupKeys))); err != nil {
		return err
	}

	getResourceResponse, err := h.resourceService.GetByKey(ctx, resource.NewGetResourceByKeyQuery(resToUpdate.GetKey()))
	if err != nil {
		return err
	}

	// retrieving groups
	groups, err := h.groupService.GetGroupsByKeys(ctx, getResourceResponse.Sharings.GetAllGroupKeys())
	if err != nil {
		l.Error("could not get groups by keys", zap.Error(err))
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, web.GetResourceResponse{
		Resource: NewResourceResponse(getResourceResponse.Resource, loggedInUser.Username, loggedInUser.Subject, groups),
	})

}
