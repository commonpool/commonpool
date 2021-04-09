package handler

import (
	"cp/pkg/api"
	"cp/pkg/memberships"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) handleGetUserGroups(c echo.Context) error {

	user, err := h.getUser(c)
	if err != nil {
		return err
	}

	var ms []*api.Membership
	if err := h.membershipStore.Find(&ms, &memberships.GetMembershipsOptions{
		UserID:  &user.ID,
		Preload: []string{"Group", "User"},
	}); err != nil {
		return err
	}

	return c.Render(http.StatusOK, "user_groups_view", map[string]interface{}{
		"Title":       user.Username + " - Groups",
		"Memberships": ms,
	})
}
