package handler

import (
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/trading"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *TradingHandler) HandleDeclineOffer(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "HandleDeclineOffer")

	offerKey, err := trading.ParseOfferKey(c.Param("id"))
	if err != nil {
		return err
	}

	err = h.tradingService.DeclineOffer(ctx, offerKey)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)

}
