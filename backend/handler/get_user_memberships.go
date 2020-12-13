package handler

import (
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetUserMemberships godoc
// @Summary Gets memberships for a given user
// @Description Gets the memberships for a given user
// @ID getUserMemberships
// @Param id path string true "ID of the user" (format:uuid)
// @Param status status MembershipStatus true "status of the membership"
// @Tags groups
// @Accept json
// @Produce json
// @Success 200 {object} web.GetUserMembershipsResponse
// @Failure 400 {object} utils.Error
// @Router /users/:id/memberships [get]
func (h *Handler) GetUserMemberships(c echo.Context) error {

	ctx, _ := GetEchoContext(c, "GetUserMemberships")

	var membershipStatus = group.AnyMembershipStatus()
	statusStr := c.QueryParam("status")
	if statusStr != "" {
		ms, err := group.ParseMembershipStatus(statusStr)
		if err != nil {
			return err
		}
		membershipStatus = &ms
	}

	memberships, err := h.groupService.GetUserMemberships(ctx, group.NewGetMembershipsForUserRequest(model.NewUserKey(c.Param("id")), membershipStatus))
	if err != nil {
		return err
	}

	groupNames, err := h.getGroupNamesForMemberships(ctx, memberships.Memberships)
	if err != nil {
		return err
	}

	userNames, err := h.getUserNamesForMemberships(ctx, memberships.Memberships)
	if err != nil {
		return err
	}

	response := web.NewGetUserMembershipsResponse(memberships.Memberships, groupNames, userNames)
	return c.JSON(http.StatusOK, response)

}
