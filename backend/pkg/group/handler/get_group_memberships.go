package handler

import (
	group2 "github.com/commonpool/backend/pkg/group"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetGroup godoc
// @Summary Gets a group memberships
// @Description Gets the members of a group
// @ID getGroupMemberships
// @Tags groups
// @Param id path string true "ID of the group" (format:uuid)
// @Param status status MembershipStatus true "status of the membership"
// @Accept json
// @Produce json
// @Success 200 {object} web.GetGroupMembershipsResponse
// @Failure 400 {object} utils.Error
// @Router /groups/:id/memberships [get]
func (h *GroupHandler) GetGroupMemberships(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "GetGroupMemberships")

	var membershipStatus = groupmodel.AnyMembershipStatus()
	statusStr := c.QueryParam("status")
	if statusStr != "" {
		ms, err := groupmodel.ParseMembershipStatus(statusStr)
		if err != nil {
			return err
		}
		membershipStatus = &ms
	}

	groupKey, err := groupmodel.ParseGroupKey(c.Param("id"))
	if err != nil {
		return err
	}

	_, err = h.groupService.GetGroup(ctx, group2.NewGetGroupRequest(groupKey))
	if err != nil {
		return err
	}

	getGroupMemberships, err := h.groupService.GetGroupMemberships(ctx, group2.NewGetMembershipsForGroupRequest(groupKey, membershipStatus))
	if err != nil {
		return err
	}

	userNames, err := h.getUserNamesForMemberships(ctx, getGroupMemberships.Memberships)
	if err != nil {
		return err
	}

	groupNames, err := h.getGroupNamesForMemberships(ctx, getGroupMemberships.Memberships)
	if err != nil {
		return err
	}

	response := web.NewGetUserMembershipsResponse(getGroupMemberships.Memberships, groupNames, userNames)
	return c.JSON(http.StatusOK, response)

}