package handler

import (
	"cp/pkg/api"
	"cp/pkg/memberships"
	"cp/pkg/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) handleGroupJoin(c echo.Context) error {

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

	if membership == nil {

		// membership does not exist. Create
		membership = &api.Membership{
			GroupID:         group.ID,
			UserID:          invitedUser.ID,
			Permission:      api.None,
			GroupConfirmed:  invitedUser.ID != authenticatedUser.ID,
			MemberConfirmed: invitedUser.ID == authenticatedUser.ID,
		}

		if err := h.membershipStore.Create(membership); err != nil {
			return err
		}

		var message string
		if invitedUser.ID != authenticatedUser.ID {
			message = fmt.Sprintf(`Successfully invited user %s to group %s. Waiting for user confirmation...`, invitedUser.HTMLLink(), group.HTMLLink())
		} else {
			message = fmt.Sprintf("Successfully sent join request to group %s. Waiting for group confirmation...", group.HTMLLink())
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

	} else {

		// Membership already exists
		if invitedUser.ID != authenticatedUser.ID {
			membership.GroupConfirmed = true
		} else {
			membership.MemberConfirmed = true
		}
		membership.Permission = api.Member

		if err := h.membershipStore.Update(membership); err != nil {
			return err
		}

		var message string
		if invitedUser.ID != authenticatedUser.ID {
			message = fmt.Sprintf("Successfully accepted user %s into group %s!", invitedUser.HTMLLink(), group.HTMLLink())
		} else {
			message = fmt.Sprintf("Successfully joined group %s!", group.HTMLLink())
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

}
