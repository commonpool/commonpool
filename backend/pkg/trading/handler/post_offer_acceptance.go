package handler

import (
	"github.com/commonpool/backend/pkg/handler"
	tradingmodel "github.com/commonpool/backend/pkg/trading/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *TradingHandler) HandleAcceptOffer(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "HandleAcceptOffer")

	offerKey, err := tradingmodel.ParseOfferKey(c.Param("id"))
	if err != nil {
		return err
	}

	err = h.tradingService.AcceptOffer(ctx, offerKey)

	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "")

}
