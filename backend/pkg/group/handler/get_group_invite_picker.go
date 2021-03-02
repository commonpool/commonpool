package handler

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/user"
	"github.com/commonpool/backend/pkg/user/handler"
	"github.com/commonpool/backend/pkg/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GetUsersForGroupInvitePickerResponse struct {
	Users []handler.UserInfoResponse `json:"users"`
	Take  int                        `json:"take"`
	Skip  int                        `json:"skip"`
}

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
		return err
	}

	take, err := utils.ParseTake(c, 10, 100)
	if err != nil {
		return err
	}

	qry := c.QueryParam("query")

	groupKey, err := keys.ParseGroupKey(c.Param("id"))
	if err != nil {
		return err
	}

	userQuery := user.Query{
		Query:      qry,
		Skip:       skip,
		Take:       take,
		NotInGroup: &groupKey,
	}

	users, err := h.userService.Find(userQuery)
	if err != nil {
		return err
	}

	responseItems := make([]handler.UserInfoResponse, len(users.Items))
	for i, u := range users.Items {
		responseItems[i] = handler.UserInfoResponse{
			Id:       u.ID,
			Username: u.Username,
		}
	}

	response := GetUsersForGroupInvitePickerResponse{
		Users: responseItems,
		Take:  take,
		Skip:  skip,
	}

	return c.JSON(http.StatusOK, response)

}
