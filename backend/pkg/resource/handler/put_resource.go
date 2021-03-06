package handler

import (
	"fmt"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	resource "github.com/commonpool/backend/pkg/resource"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type UpdateResourceRequest struct {
	Resource UpdateResourcePayload `json:"resource"`
}

type UpdateResourcePayload struct {
	Summary          string                 `json:"summary" validate:"required,max=100"`
	Description      string                 `json:"description" validate:"required,max=2000"`
	ValueInHoursFrom int                    `json:"valueInHoursFrom" validate:"min=0"`
	ValueInHoursTo   int                    `json:"valueInHoursTo" validate:"min=0"`
	SharedWith       []InputResourceSharing `json:"sharedWith"`
}

type UpdateResourceResponse struct {
	Resource Resource `json:"resource"`
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

	c.Logger().Debug("UpdateResource: updating resource")
	req := UpdateResourceRequest{}

	// Binds the request payload to the web.UpdateResourceRequest instance
	if err := c.Bind(&req); err != nil {
		return err
	}

	// Validates the UpdateResourceRequest
	if err := c.Validate(req); err != nil {
		return err
	}

	// Gets the resource id
	resourceKey, err := keys.ParseResourceKey(c.Param("id"))
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

	loggedInUser, err := oidc.GetLoggedInUser(ctx)
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
	var groupKeys []keys.GroupKey
	for _, sharing := range req.Resource.SharedWith {
		groupKey, err := keys.ParseGroupKey(sharing.GroupID)
		if err != nil {
			message := "UpdateResource: could not parse groupKey"
			c.Logger().Error(err, message)
			return c.String(http.StatusInternalServerError, message)
		}
		groupKeys = append(groupKeys, groupKey)
	}

	if err := h.resourceService.Update(ctx, resource.NewUpdateResourceQuery(resToUpdate, keys.NewGroupKeys(groupKeys))); err != nil {
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

	return c.JSON(http.StatusOK, GetResourceResponse{
		Resource: NewResourceResponse(getResourceResponse.Resource, loggedInUser.Username, loggedInUser.Subject, groups),
	})

}
