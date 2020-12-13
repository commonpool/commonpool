package handler

import (
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

// CreateOrAcceptMembership godoc
// @Summary Accept a group invitation
// @Description Accept a group invitation
// @ID acceptInvitation
// @Tags groups
// @Accept json
// @Produce json
// @Success 200 {object} web.CreateOrAcceptInvitationResponse
// @Failure 400 {object} utils.Error
// @Router /groups/memberships [post]
func (h *Handler) CreateOrAcceptMembership(c echo.Context) error {

	ctx, _ := GetEchoContext(c, "CreateOrAcceptInvitation")

	req := web.CreateOrAcceptInvitationRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	groupKey, err := group.ParseGroupKey(req.GroupID)
	if err != nil {
		return err
	}
	userKey := model.NewUserKey(req.UserID)

	membershipKey := model.NewMembershipKey(groupKey, userKey)
	acceptInvitationResponse, err := h.groupService.CreateOrAcceptInvitation(ctx, group.NewAcceptInvitationRequest(membershipKey))
	if err != nil {
		return err
	}

	memberships := group.NewMemberships([]*group.Membership{acceptInvitationResponse.Membership})

	userNames, err := h.getUserNamesForMemberships(ctx, memberships)
	if err != nil {
		return err
	}

	groupNames, err := h.getGroupNamesForMemberships(ctx, memberships)
	if err != nil {
		return err
	}

	response := web.NewCreateOrAcceptInvitationResponse(acceptInvitationResponse.Membership, groupNames, userNames)
	return c.JSON(http.StatusOK, response)

}
