package handler

import (
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/group/readmodels"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GetUserMembershipsResponse struct {
	Memberships []Membership `json:"memberships"`
}

func NewGetUserMembershipsResponse(memberships []*readmodels.MembershipReadModel) GetUserMembershipsResponse {
	responseMemberships := make([]Membership, len(memberships))
	for i, membership := range memberships {
		responseMemberships[i] = NewMembership(membership)
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

	var membershipStatus = domain.AnyMembershipStatus()
	statusStr := c.QueryParam("status")
	if statusStr != "" {
		ms, err := domain.ParseMembershipStatus(statusStr)
		if err != nil {
			return err
		}
		membershipStatus = &ms
	}

	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	userKey := loggedInUser.GetUserKey()

	userIdStr := c.Param("user_id")
	if userIdStr != "" {
		userKey = keys.NewUserKey(userIdStr)
	}

	m, err := h.getUserMemberships.Get(ctx, userKey, membershipStatus)
	if err != nil {
		return err
	}

	response := NewGetUserMembershipsResponse(m)
	return c.JSON(http.StatusOK, response)

}
