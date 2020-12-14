package handler

import (
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/handler"
	trading2 "github.com/commonpool/backend/pkg/trading"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) HandleAcceptOffer(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "HandleAcceptOffer")

	offerKey, err := model.ParseOfferKey(c.Param("id"))
	if err != nil {
		return err
	}

	_, err = h.tradingService.AcceptOffer(ctx, trading2.NewAcceptOffer(offerKey))

	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "")

}
