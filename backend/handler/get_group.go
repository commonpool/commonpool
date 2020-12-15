package handler

import (
	group2 "github.com/commonpool/backend/pkg/group"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

// GetGroup godoc
// @Summary Gets a group
// @Description Gets a group
// @ID getGroup
// @Tags groups
// @Param id path string true "ID of the group" (format:uuid)
// @Accept json
// @Produce json
// @Success 200 {object} web.GetGroupResponse
// @Failure 400 {object} utils.Error
// @Router /groups/:id [get]
func (h *Handler) GetGroup(c echo.Context) error {

	ctx, l := handler.GetEchoContext(c, "GetGroup")

	l.Debug("getting group")

	groupKey, err := groupmodel.ParseGroupKey(c.Param("id"))
	if err != nil {
		l.Error("could not parse group key", zap.Error(err))
		return NewErrResponse(c, err)
	}

	getGroupResponse, err := h.groupService.GetGroup(ctx, group2.NewGetGroupRequest(groupKey))
	if err != nil {
		l.Error("could not get group", zap.Error(err))
		return err
	}

	var response = web.NewGetGroupResponse(getGroupResponse.Group)
	return c.JSON(http.StatusOK, response)

}
