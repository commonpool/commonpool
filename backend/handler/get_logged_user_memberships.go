package handler

import (
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetLoggedInUserMemberships godoc
// @Summary Gets currently logged in user memberships
// @Description Gets the memberships for the currently logged in user
// @ID getLoggedInUserMemberships
// @Tags groups
// @Accept json
// @Produce json
// @Success 200 {object} web.GetUserMembershipsResponse
// @Failure 400 {object} utils.Error
// @Router /my/memberships [get]
func (h *Handler) GetLoggedInUserMemberships(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "GetLoggedInUserMemberships")

	authUser, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	authUserKey := authUser.GetUserKey()

	userMemberships, err := h.groupService.GetUserMemberships(ctx, group.NewGetMembershipsForUserRequest(authUserKey, group.AnyMembershipStatus()))
	if err != nil {
		return err
	}

	groupNames, err := h.getGroupNamesForMemberships(ctx, userMemberships.Memberships)
	if err != nil {
		return err
	}

	userNames, err := h.getUserNamesForMemberships(ctx, userMemberships.Memberships)
	if err != nil {
		return err
	}

	response := web.NewGetUserMembershipsResponse(userMemberships.Memberships, groupNames, userNames)
	return c.JSON(http.StatusOK, response)

}
