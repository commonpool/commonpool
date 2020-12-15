package handler

import (
	group2 "github.com/commonpool/backend/pkg/group"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	"github.com/commonpool/backend/pkg/handler"
	usermodel "github.com/commonpool/backend/pkg/user/model"
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
func (h *GroupHandler) CreateOrAcceptMembership(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "CreateOrAcceptInvitation")

	req := web.CreateOrAcceptInvitationRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	groupKey, err := groupmodel.ParseGroupKey(req.GroupID)
	if err != nil {
		return err
	}
	userKey := usermodel.NewUserKey(req.UserID)

	membershipKey := groupmodel.NewMembershipKey(groupKey, userKey)
	acceptInvitationResponse, err := h.groupService.CreateOrAcceptInvitation(ctx, group2.NewAcceptInvitationRequest(membershipKey))
	if err != nil {
		return err
	}

	memberships := groupmodel.NewMemberships([]*groupmodel.Membership{acceptInvitationResponse.Membership})

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