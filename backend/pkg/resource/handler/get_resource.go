package handler

import (
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	resource2 "github.com/commonpool/backend/pkg/resource"
	"github.com/labstack/echo/v4"
	"net/http"
)

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

	getResourceByKeyResponse, err := h.resourceService.GetByKey(ctx, resource2.NewGetResourceByKeyQuery(resourceKey))
	if err != nil {
		return err
	}
	res := getResourceByKeyResponse.Resource

	groups, err := h.groupService.GetGroupsByKeys(ctx, getResourceByKeyResponse.Sharings.GetAllGroupKeys())
	if err != nil {
		return err
	}

	ownerKey := res.GetOwnerKey()
	username, err := h.userService.GetUsername(ownerKey)
	if err != nil {
		return err
	}

	// return
	return c.JSON(http.StatusOK, GetResourceResponse{
		Resource: NewResourceResponse(res, username, ownerKey.String(), groups),
	})
}
