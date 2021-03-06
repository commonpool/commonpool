package handler

import (
	"github.com/commonpool/backend/pkg/group/readmodels"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GetGroupResponse struct {
	Group *Group `json:"group"`
}

func NewGetGroupResponse(group *readmodels.GroupReadModel) GetGroupResponse {
	return GetGroupResponse{
		Group: NewGroup(group),
	}
}

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

	ctx, _ := handler.GetEchoContext(c, "GetGroup")

	groupKey, err := keys.ParseGroupKey(c.Param("id"))
	if err != nil {
		return err
	}

	getGroupResponse, err := h.groupService.GetGroup(ctx, groupKey)
	if err != nil {
		return err
	}

	var response = NewGetGroupResponse(getGroupResponse)
	return c.JSON(http.StatusOK, response)

}
