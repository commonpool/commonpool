package handler

import (
	"cp/pkg/api"
	"cp/pkg/memberships"
	"cp/pkg/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type LeaveGroup struct {
	GroupID string `param:"groupId"`
	UserID  string `form:"userId"`
}

func (h *Handler) handleGroupLeave(c echo.Context) error {

	authenticatedUser, err := h.getAuthenticatedUser(c)
	if err != nil {
		return err
	}

	invitedUser, err := h.getUser(c)
	if err != nil {
		return err
	}

	group, err := h.getGroup(c)
	if err != nil {
		return err
	}

	membership, err := h.getMembership(c)
	if err != nil {
		return err
	}

	// If user being invited is not the currently logged in user,
	// make sure that the currently logged in user is an admin of
	// the group
	if authenticatedUser.ID != invitedUser.ID {
		var ms []*api.Membership
		requiredPermission := api.Admin
		if err := h.membershipStore.Find(&ms, &memberships.GetMembershipsOptions{
			HasPermission: &requiredPermission,
			GroupID:       &group.ID,
			UserID:        &authenticatedUser.ID,
		}); err != nil {
			return err
		}
		if len(ms) == 0 {
			return echo.ErrUnauthorized
		}
	}

	if err := h.membershipStore.Delete(membership); err != nil {
		return err
	}

	var message string
	if invitedUser.ID != authenticatedUser.ID {
		message = fmt.Sprintf("Successfully kicked user %s out from group %s", invitedUser.HTMLLink(), membership.Group.HTMLLink())
	} else {
		message = fmt.Sprintf("Successfully left group %s", membership.Group.HTMLLink())
	}

	if err := h.alertManager.AddAlert(c.Request(), c.Response().Writer, utils.Alert{
		Class:   "alert-success",
		Message: message,
	}); err != nil {
		return err
	}

	c.Response().Header().Set("Location", c.Request().Header.Get("Referer"))
	c.Response().WriteHeader(http.StatusSeeOther)
	return nil

}
