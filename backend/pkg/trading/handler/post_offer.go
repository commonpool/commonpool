package handler

import (
	"fmt"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading"
	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
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

	ctx, l := handler.GetEchoContext(c, "HandleSendOffer")

	var err error

	l.Debug("binding request to SendOfferRequest")
	req := SendOfferRequest{}
	if err = c.Bind(&req); err != nil {
		return fmt.Errorf("could not bind SendOfferRequest: %v", err)
	}

	l.Debug("validating request")
	if err = c.Validate(req); err != nil {
		return fmt.Errorf("error validating SendOfferRequest: %v", err)
	}

	l.Debug("mapping offer items")
	var tradingOfferItems []trading.OfferItem
	for i, tradingOfferItem := range req.Offer.Items {

		l.Debug("generating new key for offerItem")
		itemKey := keys.NewOfferItemKey(uuid.NewV4())

		l.Debug("mapping offer item", zap.Int("i", i))
		tradingOfferItem, err := mapNewOfferItem(tradingOfferItem, itemKey)

		if err != nil {
			return fmt.Errorf("could not map offer item %d: %v", i, err)
		}

		tradingOfferItems = append(tradingOfferItems, tradingOfferItem)
	}

	l.Debug("parsing group key")
	groupKey, err := keys.ParseGroupKey(req.Offer.GroupID)
	if err != nil {
		return err
	}

	l.Debug("sending offer")
	offer, offerItems, err := h.tradingService.SendOffer(ctx, groupKey, trading.NewOfferItems(tradingOfferItems), "")
	if err != nil {
		return err
	}

	l.Debug("getting approvers for offer")
	approvers, err := h.tradingService.FindApproversForOffer(offer.Key)
	if err != nil {
		return err
	}

	l.Debug("mapping to web offer")
	webOffer, err := h.mapToWebOffer(offer, offerItems, approvers)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, &OfferResponse{
		Offer: webOffer,
	})

}
