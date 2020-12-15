package handler

import (
	"github.com/commonpool/backend/pkg/handler"
	tradingmodel "github.com/commonpool/backend/pkg/trading/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *TradingHandler) HandleConfirmServiceProvided(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "HandleConfirmServiceProvided")

	offerItemKey, err := tradingmodel.ParseOfferItemKey(c.Param("id"))
	if err != nil {
		return err
	}

	err = h.tradingService.ConfirmServiceProvided(ctx, offerItemKey)
	if err != nil {
		return err
	}

	offerItem, err := h.tradingService.GetOfferItem(ctx, offerItemKey)
	if err != nil {
		return err
	}

	approvers, err := h.tradingService.FindApproversForOffer(offerItem.GetOfferKey())
	if err != nil {
		return err
	}

	webResponse, err := mapWebOfferItem(offerItem, approvers)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, webResponse)

}
