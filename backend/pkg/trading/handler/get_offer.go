package handler

import (
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *TradingHandler) HandleGetOffer(c echo.Context) error {
	ctx := handler.GetContext(c)
	offerKey, err := keys.ParseOfferKey(c.Param("id"))
	if err != nil {
		return err
	}
	offer, err := h.getOffer.Get(ctx, offerKey)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, offer)
}
