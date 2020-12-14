package handler

import (
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/exceptions"
	group2 "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"
	"net/http"
	"strings"
)

// CreateGroup godoc
// @Summary Creates a group
// @Description Creates a group and sets the authenticated user as the owner
// @ID createGroup
// @Tags groups
// @Param group body web.CreateGroupRequest true "Group to create"
// @Accept json
// @Produce json
// @Success 200 {object} web.CreateGroupResponse
// @Failure 400 {object} utils.Error
// @Router /groups [post]
func (h *Handler) CreateGroup(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "CreateGroup")

	req := web.CreateGroupRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

	if req.Name == "" {
		return NewErrResponse(c, exceptions.ErrValidation("name is required"))
	}

	var groupKey = model.NewGroupKey(uuid.NewV4())

	createGroupResponse, err := h.groupService.CreateGroup(ctx, group2.NewCreateGroupRequest(groupKey, req.Name, req.Description))
	if err != nil {
		return err
	}

	var response = web.NewCreateGroupResponse(createGroupResponse.Group)
	return c.JSON(http.StatusCreated, response)

}
