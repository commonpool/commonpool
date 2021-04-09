package handler

import (
	"cp/pkg/api"
	"cp/pkg/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type SetPermission struct {
	Permission api.MembershipPermission `form:"permission"`
}

func (h *Handler) handleGroupSetPermission(c echo.Context) error {

	user, err := h.getUser(c)
	if err != nil {
		return err
	}

	membership, err := h.getMembership(c)
	if err != nil {
		return err
	}

	authenticatedUserMembership, err := h.getAuthenticatedUserMembership(c)
	if err != nil {
		return err
	}

	// Make sure the authenticated user is an admin
	if !authenticatedUserMembership.IsActive() || !authenticatedUserMembership.IsAdmin() {
		return echo.ErrForbidden
	}

	// Make sure the authenticated user has greater permissions
	// than the changed user permission
	if !authenticatedUserMembership.Permission.Gte(membership.Permission) {
		return echo.ErrForbidden
	}

	var payload SetPermission
	if err := c.Bind(&payload); err != nil {
		return err
	}

	// make sure the new permissions are lesser or equal than the
	// authenticated user permission
	if !authenticatedUserMembership.Permission.Gte(payload.Permission) {
		return echo.ErrForbidden
	}

	// Make sure it's a valid permission
	if payload.Permission != api.Owner && payload.Permission != api.Admin && payload.Permission != api.Member {
		return echo.ErrBadRequest
	}

	membership.Permission = payload.Permission
	if err := h.membershipStore.Update(membership); err != nil {
		return err
	}

	if err := h.alertManager.AddAlert(c.Request(), c.Response().Writer, utils.Alert{
		Class:   "alert-success",
		Message: fmt.Sprintf("Successfully assigned <b>%s</b> permissions to user %s", payload.Permission, user.HTMLLink()),
	}); err != nil {
		return err
	}

	c.Response().Header().Set("Location", c.Request().Header.Get("Referer"))
	c.Response().WriteHeader(http.StatusSeeOther)

	return nil

}
