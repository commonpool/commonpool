package handler

import (
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/handler"
	trading2 "github.com/commonpool/backend/pkg/trading"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"
	"net/http"
)

func (h *Handler) HandleSendOffer(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "HandleSendOffer")

	var err error

	req := web.SendOfferRequest{}
	if err = c.Bind(&req); err != nil {
		return err
	}

	if err = c.Validate(req); err != nil {
		return err
	}

	var tradingOfferItems []trading2.OfferItem
	for _, tradingOfferItem := range req.Offer.Items {
		itemKey := model.NewOfferItemKey(uuid.NewV4())
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

	offer, offerItems, err := h.tradingService.SendOffer(ctx, groupKey, trading2.NewOfferItems(tradingOfferItems), "")
	if err != nil {
		return err
	}

	approvers, err := h.tradingStore.FindApproversForOffer(offer.Key)
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
