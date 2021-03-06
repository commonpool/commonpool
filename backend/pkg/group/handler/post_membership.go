package handler

import (
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/group/readmodels"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CreateOrAcceptInvitationRequest struct {
	UserID  string `json:"userId"`
	GroupID string `json:"groupId"`
}

type CreateOrAcceptInvitationResponse struct {
	Membership Membership `json:"membership"`
}

func NewCreateOrAcceptInvitationResponse(membership *readmodels.MembershipReadModel) *CreateOrAcceptInvitationResponse {
	return &CreateOrAcceptInvitationResponse{
		Membership: NewMembership(membership),
	}
}

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

	ctx, _ := handler.GetEchoContext(c, "CreateOrAcceptInvitation")

	req := CreateOrAcceptInvitationRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	groupKey, err := keys.ParseGroupKey(req.GroupID)
	if err != nil {
		return err
	}
	userKey := keys.NewUserKey(req.UserID)

	membershipKey := keys.NewMembershipKey(groupKey, userKey)
	err = h.groupService.CreateOrAcceptInvitation(ctx, group.NewAcceptInvitationRequest(membershipKey))
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)

	// memberships := domain.NewMemberships([]*domain.Membership{acceptInvitationResponse.Membership})
	//
	// userNames, err := h.getUserNamesForMemberships(ctx, memberships)
	// if err != nil {
	// 	return err
	// }
	//
	// groupNames, err := h.getGroupNamesForMemberships(ctx, memberships)
	// if err != nil {
	// 	return err
	// }
	//
	// response := NewCreateOrAcceptInvitationResponse(acceptInvitationResponse.Membership, groupNames, userNames)
	// return c.JSON(http.StatusOK, response)

}
