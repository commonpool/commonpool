package handler

import (
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetMembership godoc
// @Summary Gets the membership for a given user and group
// @Description Gets the membership for a given user and group
// @ID getMembership
// @Param groupId path string true "ID of the group" (format:uuid)
// @Param userId path string true "ID of the user"
// @Tags groups
// @Accept json
// @Produce json
// @Success 200 {object} web.GetMembershipResponse
// @Failure 400 {object} utils.Error
// @Router /groups/:groupId/memberships/:userId [get]
func (h *Handler) GetMembership(c echo.Context) error {

	ctx, _ := GetEchoContext(c, "GetMembership")

	userKey := model.NewUserKey(c.Param("userId"))

	groupKey, err := group.ParseGroupKey(c.Param("groupId"))
	if err != nil {
		return err
	}

	getMemberships, err := h.groupService.GetMembership(ctx, group.NewGetMembershipRequest(model.NewMembershipKey(groupKey, userKey)))
	if err != nil {
		return err
	}

	var memberships = group.NewMemberships([]*group.Membership{getMemberships.Membership})

	groupNames, err := h.getGroupNamesForMemberships(ctx, memberships)
	if err != nil {
		return err
	}

	userNames, err := h.getUserNamesForMemberships(ctx, memberships)
	if err != nil {
		return err
	}

	response := web.NewGetMembershipResponse(getMemberships.Membership, groupNames, userNames)
	return c.JSON(http.StatusOK, response)

}
