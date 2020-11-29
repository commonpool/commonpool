package handler

import (
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/trading"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) Chatback(c echo.Context) error {

	// Unmarshal request
	req := web.InteractionCallback{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	c.Logger().Info(req.Payload.Actions[0].ActionID)

	if req.Payload.Actions[0].ActionID == "accept_offer" {
		return h.HandleChatbackOfferAccepted(c, req)
	} else if req.Payload.Actions[0].ActionID == "confirm_item_received" {
		return h.HandleChatbackConfirmItemReceived(c, req)
	} else if req.Payload.Actions[0].ActionID == "confirm_item_given" {
		return h.HandleChatbackConfirmItemGiven(c, req)
	}

	return nil

}

func (h *Handler) HandleChatbackConfirmItemReceived(c echo.Context, req web.InteractionCallback) error {
	return h.HandleChatbackConfirmItemGivenOrReceived(
		c,
		req,
		trading.OfferItemReceiving)
}

func (h *Handler) HandleChatbackConfirmItemGiven(c echo.Context, req web.InteractionCallback) error {
	return h.HandleChatbackConfirmItemGivenOrReceived(
		c,
		req,
		trading.OfferItemGiving)

}

func (h *Handler) HandleChatbackConfirmItemGivenOrReceived(
	c echo.Context,
	req web.InteractionCallback,
	bond trading.OfferItemBond) error {

	var err error

	ctx, l := GetEchoContext(c, "HandleChatbackConfirmItemGivenOrReceived")

	l.Debug("retrieving offer item id from payload")

	// retrieving item id from payload
	itemId := req.Payload.Actions[0].Value
	if itemId == nil {
		l.Error("value is required")
		return c.String(http.StatusBadRequest, "value is required")
	}

	l.Debug("converting item id to item key")

	// converting item id to item key
	offerItemKey, err := model.ParseOfferItemKey(*itemId)
	if err != nil {
		l.Error("could not get offer item id from request", zap.Error(err))
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = h.tradingService.ConfirmItemReceivedOrGiven(ctx, offerItemKey)

	if err != nil {
		l.Error("could not confirm item was received or given", zap.Error(err))
		return err
	}

	return c.JSON(http.StatusOK, "OK")

}

func (h *Handler) HandleChatbackOfferAccepted(c echo.Context, req web.InteractionCallback) error {

	ctx, l := GetEchoContext(c, "HandleChatbackOfferAccepted")

	offerId := req.Payload.Actions[0].Value
	if offerId == nil {
		l.Error("offerId value is required")
		return c.String(http.StatusBadRequest, "value is required")
	}

	uid, err := uuid.FromString(*offerId)
	if err != nil {
		l.Error("could not parse offer id")
		return c.String(http.StatusBadRequest, err.Error())
	}

	offerKey := model.NewOfferKey(uid)

	_, err = h.tradingService.AcceptOffer(ctx, trading.NewAcceptOffer(offerKey))
	if err != nil {
		l.Error("could not accept offer", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "OK")
}
