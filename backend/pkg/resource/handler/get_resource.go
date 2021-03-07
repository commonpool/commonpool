package handler

import (
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
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
	g, ctx := errgroup.WithContext(ctx)

	resourceKey, err := keys.ParseResourceKey(c.Param("id"))
	if err != nil {
		return err
	}

	var resource *readmodel.ResourceReadModel
	var shares []*readmodel.ResourceSharingReadModel

	g.Go(func() error {
		resource, err = h.getResource.Get(ctx, resourceKey)
		return err
	})
	g.Go(func() error {
		shares, err = h.getResourceSharings.Get(ctx, resourceKey)
		return err
	})
	if err := g.Wait(); err != nil {
		return err
	}

	// return
	return c.JSON(http.StatusOK, GetResourceResponse{
		Resource: NewResourceResponse(resource, shares),
	})
}
