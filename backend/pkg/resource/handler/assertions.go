package handler

import (
	group2 "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

func (h *ResourceHandler) ensureResourceIsSharedWithGroupsTheUserIsActiveMemberOf(c echo.Context, loggedInUserKey keys.UserKey, sharedWithGroups *keys.GroupKeys) (error, bool) {

	ctx, l := handler.GetEchoContext(c, "ensureResourceIsSharedWithGroupsTheUserIsActiveMemberOf")

	var membershipStatus = group2.ApprovedMembershipStatus

	userMemberships, err := h.groupService.GetUserMemberships(ctx, group2.NewGetMembershipsForUserRequest(loggedInUserKey, &membershipStatus))
	if err != nil {
		l.Error("could not get user memberships", zap.Error(err))
		return err, true
	}

	// Checking if resource is shared with groups the user is part of
	for _, sharedWith := range sharedWithGroups.Items {
		hasMembershipInGroup := userMemberships.Memberships.ContainsMembershipForGroup(sharedWith)
		if !hasMembershipInGroup {
			return c.String(http.StatusBadRequest, "cannot share resource with a group you are not part of"), true
		}
	}
	return nil, false
}
