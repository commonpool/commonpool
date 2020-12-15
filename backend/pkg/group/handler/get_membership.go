package handler

import (
	group2 "github.com/commonpool/backend/pkg/group"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	"github.com/commonpool/backend/pkg/handler"
	usermodel "github.com/commonpool/backend/pkg/user/model"
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
func (h *GroupHandler) GetMembership(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "GetMembership")

	userKey := usermodel.NewUserKey(c.Param("userId"))

	groupKey, err := groupmodel.ParseGroupKey(c.Param("id"))
	if err != nil {
		return err
	}

	getMemberships, err := h.groupService.GetMembership(ctx, group2.NewGetMembershipRequest(groupmodel.NewMembershipKey(groupKey, userKey)))
	if err != nil {
		return err
	}

	var memberships = groupmodel.NewMemberships([]*groupmodel.Membership{getMemberships.Membership})

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
