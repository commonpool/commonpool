package handler

import (
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) HandleConfirmBorrowedResourceReturned(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "HandleConfirmBorrowedResourceReturned")

	offerItemKey, err := model.ParseOfferItemKey(c.Param("id"))
	if err != nil {
		return err
	}

	err = h.tradingService.ConfirmBorrowedResourceReturned(ctx, offerItemKey)
	if err != nil {
		return err
	}

	offerItem, err := h.tradingService.GetOfferItem(ctx, offerItemKey)
	if err != nil {
		return err
	}

	approvers, err := h.tradingStore.FindApproversForOffer(offerItem.GetOfferKey())
	if err != nil {
		return err
	}

	webResponse, err := mapWebOfferItem(offerItem, approvers)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, webResponse)

}
