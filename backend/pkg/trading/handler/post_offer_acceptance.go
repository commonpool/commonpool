package handler

import (
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *TradingHandler) HandleAcceptOffer(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "HandleAcceptOffer")
	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}

	offerKey, err := keys.ParseOfferKey(c.Param("id"))
	if err != nil {
		return err
	}

	offer, err := h.offerRepo.Load(ctx, offerKey)
	if err != nil {
		return err
	}

	if offer.GetVersion() == 0 {
		return exceptions.ErrOfferNotFound
	}

	permission, err := h.getPermissions.Get(ctx, offerKey)
	if err != nil {
		return err
	}

	if err := offer.ApproveAll(loggedInUser.GetUserKey(), permission); err != nil {
		return err
	}

	if err := h.offerRepo.Save(ctx, offer); err != nil {
		return err
	}

	return c.String(http.StatusOK, "")

}
