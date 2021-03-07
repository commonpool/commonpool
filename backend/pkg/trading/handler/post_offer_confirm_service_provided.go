package handler

import (
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *TradingHandler) HandleConfirmServiceProvided(c echo.Context) error {
	ctx, _ := handler.GetEchoContext(c, "HandleConfirmServiceProvided")
	offerItemKey, err := keys.ParseOfferItemKey(c.Param("id"))
	if err != nil {
		return err
	}
	offerItem, err := h.doWithOfferItem(ctx, offerItemKey, func(offer *domain.Offer) error {
		loggedInUser, err := oidc.GetLoggedInUser(ctx)
		if err != nil {
			return err
		}
		return offer.NotifyServiceGiven(loggedInUser.GetUserKey(), offerItemKey)
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, offerItem)
}
