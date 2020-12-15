package handler

import (
	"github.com/commonpool/backend/pkg/handler"
	tradingmodel "github.com/commonpool/backend/pkg/trading/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) HandleDeclineOffer(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "HandleDeclineOffer")

	offerKey, err := tradingmodel.ParseOfferKey(c.Param("id"))
	if err != nil {
		return err
	}

	err = h.tradingService.DeclineOffer(ctx, offerKey)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)

}
