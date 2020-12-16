package handler

import (
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GetGroupResponse struct {
	Group *Group `json:"group"`
}

func NewGetGroupResponse(group *group.Group) GetGroupResponse {
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

	groupKey, err := group.ParseGroupKey(c.Param("id"))
	if err != nil {
		return err
	}

	getGroupResponse, err := h.groupService.GetGroup(ctx, group.NewGetGroupRequest(groupKey))
	if err != nil {
		return err
	}

	var response = NewGetGroupResponse(getGroupResponse.Group)
	return c.JSON(http.StatusOK, response)

}
