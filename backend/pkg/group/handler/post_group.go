package handler

import (
	retry "github.com/avast/retry-go"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/group/readmodels"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type CreateGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateGroupResponse struct {
	Group *Group `json:"group"`
}

func NewCreateGroupResponse(group *readmodels.GroupReadModel) CreateGroupResponse {
	return CreateGroupResponse{
		Group: NewGroup(group),
	}
}

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

	req := CreateGroupRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)

	if req.Name == "" {
		return exceptions.ErrValidation("name is required")
	}

	groupKey, err := h.groupService.CreateGroup(ctx, group.NewCreateGroupRequest(req.Name, req.Description))
	if err != nil {
		return err
	}

	errChan := make(chan error)
	grpChan := make(chan *readmodels.GroupReadModel)

	go func() {
		err = retry.Do(func() error {
			grp, err := h.groupService.GetGroup(ctx, groupKey)
			if err != nil {
				return err
			}
			grpChan <- grp
			return nil
		})
		errChan <- err
	}()

	select {
	case err := <-errChan:
		return err
	case grp := <-grpChan:
		var response = NewCreateGroupResponse(grp)
		return c.JSON(http.StatusCreated, response)
	}

}
