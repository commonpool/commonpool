package handler

import (
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/user"
	"github.com/commonpool/backend/utils"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetGroup godoc
// @Summary User picker for group invite
// @Description Finds users to invite on a group
// @ID inviteMemberPicker
// @Tags groups
// @Param id path string true "ID of the group" (format:uuid)
// @Accept json
// @Produce json
// @Success 200 {object} web.GetGroupMembershipsResponse
// @Failure 400 {object} utils.Error
// @Router /groups/:id/invite-member-picker [get]
func (h *Handler) GetUsersForGroupInvitePicker(c echo.Context) error {
	skip, err := utils.ParseSkip(c)
	if err != nil {
		return NewErrResponse(c, err)
	}

	take, err := utils.ParseTake(c, 10, 100)
	if err != nil {
		return NewErrResponse(c, err)
	}

	qry := c.QueryParam("query")

	groupKey, err := model.ParseGroupKey(c.Param("id"))
	if err != nil {
		return NewErrResponse(c, err)
	}

	userQuery := user.Query{
		Query:      qry,
		Skip:       skip,
		Take:       take,
		NotInGroup: &groupKey,
	}

	users, err := h.authStore.Find(userQuery)
	if err != nil {
		return NewErrResponse(c, err)
	}

	responseItems := make([]web.UserInfoResponse, len(users))
	for i, u := range users {
		responseItems[i] = web.UserInfoResponse{
			Id:       u.ID,
			Username: u.Username,
		}
	}

	response := web.GetUsersForGroupInvitePickerResponse{
		Users: responseItems,
		Take:  take,
		Skip:  skip,
	}

	return c.JSON(http.StatusOK, response)

}
