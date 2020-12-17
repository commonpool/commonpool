package handler

import (
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading"
	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"
	"net/http"
)

type SendOfferRequest struct {
	Offer SendOfferPayload `json:"offer" validate:"required"`
}

type SendOfferPayload struct {
	Items   []SendOfferPayloadItem `json:"items" validate:"min=1"`
	GroupID string                 `json:"groupId" validate:"uuid"`
	Message string                 `json:"message"`
}

type SendOfferPayloadItem struct {
	Type       trading.OfferItemType `json:"type"`
	To         OfferItemTarget       `json:"to" validate:"required,uuid"`
	From       *OfferItemTarget      `json:"from" validate:"required,uuid"`
	ResourceId *string               `json:"resourceId" validate:"required,uuid"`
	Duration   *string               `json:"duration"`
	Amount     *string               `json:"amount"`
}

func (h *TradingHandler) HandleSendOffer(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "HandleSendOffer")

	var err error

	req := SendOfferRequest{}
	if err = c.Bind(&req); err != nil {
		return err
	}

	if err = c.Validate(req); err != nil {
		return err
	}

	var tradingOfferItems []trading.OfferItem
	for _, tradingOfferItem := range req.Offer.Items {
		itemKey := trading.NewOfferItemKey(uuid.NewV4())
		tradingOfferItem, err := mapNewOfferItem(tradingOfferItem, itemKey)
		if err != nil {
			return err
		}
		tradingOfferItems = append(tradingOfferItems, tradingOfferItem)
	}

	groupKey, err := keys.ParseGroupKey(req.Offer.GroupID)
	if err != nil {
		return err
	}

	offer, offerItems, err := h.tradingService.SendOffer(ctx, groupKey, trading.NewOfferItems(tradingOfferItems), "")
	if err != nil {
		return err
	}

	approvers, err := h.tradingService.FindApproversForOffer(offer.Key)
	if err != nil {
		return err
	}

	webOffer, err := h.mapToWebOffer(offer, offerItems, approvers)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, &OfferResponse{
		Offer: webOffer,
	})

}
