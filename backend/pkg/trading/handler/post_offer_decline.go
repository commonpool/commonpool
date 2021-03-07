package handler

import (
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *TradingHandler) HandleDeclineOffer(c echo.Context) error {
	ctx, _ := handler.GetEchoContext(c, "HandleDeclineOffer")

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

	if err := offer.DeclineOffer(loggedInUser.GetUserKey()); err != nil {
		return err
	}

	if err := h.offerRepo.Save(ctx, offer); err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)
}
