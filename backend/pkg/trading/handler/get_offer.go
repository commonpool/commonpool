package handler

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"
	"net/http"
)

func (h *TradingHandler) HandleGetOffer(c echo.Context) error {

	var err error

	offerIdStr := c.Param("id")

	offerId, err := uuid.FromString(offerIdStr)
	if err != nil {
		return err
	}
	offerKey := keys.NewOfferKey(offerId)

	offer, err := h.getWebOffer(offerKey)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, offer)

}

func (h *TradingHandler) getWebOffer(offerKey keys.OfferKey) (*OfferResponse, error) {

	offer, err := h.tradingService.GetOffer(offerKey)
	if err != nil {
		return nil, err
	}

	items, err := h.tradingService.GetOfferItemsForOffer(offerKey)
	if err != nil {
		return nil, err
	}

	approvers, err := h.tradingService.FindApproversForOffer(offer.Key)
	if err != nil {
		return nil, err
	}

	webOffer, err := h.mapToWebOffer(offer, items, approvers)
	if err != nil {
		return nil, err
	}

	response := OfferResponse{
		Offer: webOffer,
	}

	return &response, nil
}
