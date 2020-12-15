package handler

import (
	group2 "github.com/commonpool/backend/pkg/group"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	"github.com/commonpool/backend/pkg/handler"
	usermodel "github.com/commonpool/backend/pkg/user/model"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

// CancelOrDeclineInvitation godoc
// @Summary declines a group invitation
// @Description declines a group invitation
// @ID declineInvitation
// @Tags groups
// @Accept json
// @Produce json
// @Success 202 {object} web.CancelOrDeclineInvitationResponse
// @Failure 400 {object} utils.Error
// @Router /memberships [delete]
func (h *Handler) CancelOrDeclineInvitation(c echo.Context) error {

	ctx, l := handler.GetEchoContext(c, "CancelOrDeclineInvitation")

	req := web.CancelOrDeclineInvitationRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	groupKey, err := groupmodel.ParseGroupKey(req.GroupID)
	if err != nil {
		return err
	}
	userKey := usermodel.NewUserKey(req.UserID)

	membershipKey := groupmodel.NewMembershipKey(groupKey, userKey)

	err = h.groupService.CancelOrDeclineInvitation(ctx, group2.NewDelineInvitationRequest(membershipKey))
	if err != nil {
		l.Error("could not decline invitation", zap.Error(err))
		return err
	}

	return c.NoContent(http.StatusAccepted)

}
