package handler

import (
	"github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/group/readmodels"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GetGroupMembershipsResponse struct {
	Memberships []Membership `json:"memberships"`
}

func NewGetGroupMembershipsResponse(memberships []*readmodels.MembershipReadModel) GetUserMembershipsResponse {
	responseMemberships := make([]Membership, len(memberships))
	for i, membership := range memberships {
		responseMemberships[i] = NewMembership(membership)
	}
	return GetUserMembershipsResponse{
		Memberships: responseMemberships,
	}
}

// GetGroup godoc
// @Summary Gets a group memberships
// @Description Gets the members of a group
// @ID getGroupMemberships
// @Tags groups
// @Param id path string true "ID of the group" (format:uuid)
// @Param status status MembershipStatus true "status of the membership"
// @Accept json
// @Produce json
// @Success 200 {object} web.GetGroupMembershipsResponse
// @Failure 400 {object} utils.Error
// @Router /groups/:id/memberships [get]
func (h *Handler) GetGroupMemberships(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "GetGroupMemberships")

	var membershipStatus = domain.AnyMembershipStatus()
	statusStr := c.QueryParam("status")
	if statusStr != "" {
		ms, err := domain.ParseMembershipStatus(statusStr)
		if err != nil {
			return err
		}
		membershipStatus = &ms
	}

	groupKey, err := keys.ParseGroupKey(c.Param("id"))
	if err != nil {
		return err
	}

	m, err := h.getGroupMemberships.Get(ctx, groupKey, membershipStatus)
	if err != nil {
		return err
	}

	response := NewGetGroupMembershipsResponse(m)
	return c.JSON(http.StatusOK, response)

}
