package handler

import (
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/handler"
	tradingmodel "github.com/commonpool/backend/pkg/trading/model"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"
	"net/http"
)

func (h *TradingHandler) HandleSendOffer(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "HandleSendOffer")

	var err error

	req := web.SendOfferRequest{}
	if err = c.Bind(&req); err != nil {
		return err
	}

	if err = c.Validate(req); err != nil {
		return err
	}

	var tradingOfferItems []tradingmodel.OfferItem
	for _, tradingOfferItem := range req.Offer.Items {
		itemKey := tradingmodel.NewOfferItemKey(uuid.NewV4())
		tradingOfferItem, err := mapNewOfferItem(tradingOfferItem, itemKey)
		if err != nil {
			return err
		}
		tradingOfferItems = append(tradingOfferItems, tradingOfferItem)
	}

	groupKey, err := group.ParseGroupKey(req.Offer.GroupID)
	if err != nil {
		return err
	}

	offer, offerItems, err := h.tradingService.SendOffer(ctx, groupKey, tradingmodel.NewOfferItems(tradingOfferItems), "")
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

	return c.JSON(http.StatusCreated, &web.GetOfferResponse{
		Offer: webOffer,
	})

}
