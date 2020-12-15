package handler

import (
	"github.com/commonpool/backend/pkg/handler"
	resource2 "github.com/commonpool/backend/pkg/resource"
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
	"github.com/commonpool/backend/web"
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
func (h *Handler) GetResource(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "GetResource")

	resourceKey, err := resourcemodel.ParseResourceKey(c.Param("id"))
	if err != nil {
		return err
	}

	getResourceByKeyResponse, err := h.resourceStore.GetByKey(ctx, resource2.NewGetResourceByKeyQuery(resourceKey))
	if err != nil {
		return err
	}
	res := getResourceByKeyResponse.Resource

	groups, err := h.groupService.GetGroupsByKeys(ctx, getResourceByKeyResponse.Sharings.GetAllGroupKeys())
	if err != nil {
		return err
	}

	ownerKey := res.GetOwnerKey()
	username, err := h.authStore.GetUsername(ownerKey)
	if err != nil {
		return err
	}

	// return
	return c.JSON(http.StatusOK, web.GetResourceResponse{
		Resource: NewResourceResponse(res, username, ownerKey.String(), groups),
	})
}
