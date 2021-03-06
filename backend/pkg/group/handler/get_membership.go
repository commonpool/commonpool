package handler

import (
	"github.com/avast/retry-go"
	"github.com/commonpool/backend/pkg/group/readmodels"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GetMembershipResponse struct {
	Membership Membership `json:"membership"`
}

func NewGetMembershipResponse(membership *readmodels.MembershipReadModel) *GetMembershipResponse {
	return &GetMembershipResponse{
		Membership: NewMembership(membership),
	}
}

// GetMembership godoc
// @Summary Gets the membership for a given user and group
// @Description Gets the membership for a given user and group
// @ID getMembership
// @Param groupId path string true "ID of the group" (format:uuid)
// @Param userId path string true "ID of the user"
// @Tags groups
// @Accept json
// @Produce json
// @Success 200 {object} web.GetMembershipResponse
// @Failure 400 {object} utils.Error
// @Router /groups/:groupId/memberships/:userId [get]
func (h *Handler) GetMembership(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "GetMembership")

	userKey := keys.NewUserKey(c.Param("userId"))

	groupKey, err := keys.ParseGroupKey(c.Param("id"))
	if err != nil {
		return err
	}

	membershipKey := keys.NewMembershipKey(groupKey, userKey)
	var membership *readmodels.MembershipReadModel
	err = retry.Do(func() error {
		var err error
		membership, err = h.getMembership.Get(ctx, membershipKey)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	resp := &GetMembershipResponse{
		Membership: Membership{
			UserID:         membership.UserKey,
			GroupID:        membership.GroupKey,
			IsAdmin:        membership.IsAdmin,
			IsMember:       membership.IsMember,
			IsOwner:        membership.IsOwner,
			GroupConfirmed: membership.GroupConfirmed,
			UserConfirmed:  membership.UserConfirmed,
			GroupName:      membership.GroupName,
			UserName:       membership.UserName,
		},
	}

	return c.JSON(http.StatusOK, resp)

}
