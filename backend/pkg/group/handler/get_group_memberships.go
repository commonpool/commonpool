package handler

import (
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GetGroupMembershipsResponse struct {
	Memberships []Membership `json:"memberships"`
}

func NewGetGroupMembershipsResponse(memberships *group.Memberships, groupNames group.Names, userNames auth.UserNames) GetUserMembershipsResponse {
	responseMemberships := make([]Membership, len(memberships.Items))
	for i, membership := range memberships.Items {
		responseMemberships[i] = NewMembership(membership, groupNames, userNames)
	}
	return GetUserMembershipsResponse{
		Memberships: responseMemberships,
	}
}

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
func (h *Handler) GetGroupMemberships(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "GetGroupMemberships")

	var membershipStatus = group.AnyMembershipStatus()
	statusStr := c.QueryParam("status")
	if statusStr != "" {
		ms, err := group.ParseMembershipStatus(statusStr)
		if err != nil {
			return err
		}
		membershipStatus = &ms
	}

	groupKey, err := keys.ParseGroupKey(c.Param("id"))
	if err != nil {
		return err
	}

	_, err = h.groupService.GetGroup(ctx, group.NewGetGroupRequest(groupKey))
	if err != nil {
		return err
	}

	getGroupMemberships, err := h.groupService.GetGroupMemberships(ctx, group.NewGetMembershipsForGroupRequest(groupKey, membershipStatus))
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

	response := NewGetGroupMembershipsResponse(getGroupMemberships.Memberships, groupNames, userNames)
	return c.JSON(http.StatusOK, response)

}
