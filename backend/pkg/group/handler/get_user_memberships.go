package handler

import (
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GetUserMembershipsResponse struct {
	Memberships []Membership `json:"memberships"`
}

func NewGetUserMembershipsResponse(memberships *group.Memberships, groupNames group.Names, userNames auth.UserNames) GetUserMembershipsResponse {
	responseMemberships := make([]Membership, len(memberships.Items))
	for i, membership := range memberships.Items {
		responseMemberships[i] = NewMembership(membership, groupNames, userNames)
	}
	return GetUserMembershipsResponse{
		Memberships: responseMemberships,
	}
}

// GetUserMemberships godoc
// @Summary Gets memberships for a given user
// @Description Gets the memberships for a given user
// @ID getUserMemberships
// @Param user_id query string false "ID of the user. If not set, defaults to the logged in user id" (format:uuid)
// @Param status status MembershipStatus true "status of the membership"
// @Tags groups
// @Accept json
// @Produce json
// @Success 200 {object} web.GetUserMembershipsResponse
// @Failure 400 {object} utils.Error
// @Router /memberships [get]
func (h *Handler) GetUserMemberships(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "GetUserMemberships")

	var membershipStatus = group.AnyMembershipStatus()
	statusStr := c.QueryParam("status")
	if statusStr != "" {
		ms, err := group.ParseMembershipStatus(statusStr)
		if err != nil {
			return err
		}
		membershipStatus = &ms
	}

	loggedInUser, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	userKey := loggedInUser.GetUserKey()

	userIdStr := c.Param("user_id")
	if userIdStr != "" {
		userKey = keys.NewUserKey(userIdStr)
	}

	memberships, err := h.groupService.GetUserMemberships(ctx, group.NewGetMembershipsForUserRequest(userKey, membershipStatus))
	if err != nil {
		return err
	}

	groupNames, err := h.getGroupNamesForMemberships(ctx, memberships.Memberships)
	if err != nil {
		return err
	}

	userNames, err := h.getUserNamesForMemberships(ctx, memberships.Memberships)
	if err != nil {
		return err
	}

	response := NewGetUserMembershipsResponse(memberships.Memberships, groupNames, userNames)
	return c.JSON(http.StatusOK, response)

}
