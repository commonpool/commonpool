package handler

import (
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/handler"
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
	} else if req.Payload.Actions[0].ActionID == "confirm_service_provided" {
		return h.HandleChatbackConfirmServiceProvided(c, req)
	} else if req.Payload.Actions[0].ActionID == "confirm_resource_transferred" {
		return h.HandleChatbackConfirmResourceTransferred(c, req)
	} else if req.Payload.Actions[0].ActionID == "confirm_resource_borrowed" {
		return h.HandleChatbackConfirmResourceBorrowed(c, req)
	} else if req.Payload.Actions[0].ActionID == "confirm_resource_borrowed_returned" {
		return h.HandleChatbackConfirmResourceBorrowedReturned(c, req)
	}

	return nil

}

func (h *Handler) HandleChatbackConfirmServiceProvided(c echo.Context, req web.InteractionCallback) error {

	ctx, l := handler.GetEchoContext(c, "HandleChatbackConfirmServiceProvided")

	// retrieving item id from payload
	offerItemId := req.Payload.Actions[0].Value
	if offerItemId == nil {
		l.Error("value is required")
		return c.String(http.StatusBadRequest, "value is required")
	}

	// converting item id to item key
	offerItemKey, err := model.ParseOfferItemKey(*offerItemId)
	if err != nil {
		l.Error("could not get offer item id from request", zap.Error(err))
		return c.String(http.StatusBadRequest, err.Error())
	}

	return h.tradingService.ConfirmServiceProvided(ctx, offerItemKey)

}

func (h *Handler) HandleChatbackConfirmResourceTransferred(c echo.Context, req web.InteractionCallback) error {

	ctx, l := handler.GetEchoContext(c, "HandleChatbackConfirmResourceTransferred")

	// retrieving item id from payload
	offerItemId := req.Payload.Actions[0].Value
	if offerItemId == nil {
		l.Error("value is required")
		return c.String(http.StatusBadRequest, "value is required")
	}

	// converting item id to item key
	offerItemKey, err := model.ParseOfferItemKey(*offerItemId)
	if err != nil {
		l.Error("could not get offer item id from request", zap.Error(err))
		return c.String(http.StatusBadRequest, err.Error())
	}

	return h.tradingService.ConfirmResourceTransferred(ctx, offerItemKey)

}

func (h *Handler) HandleChatbackConfirmResourceBorrowed(c echo.Context, req web.InteractionCallback) error {

	ctx, l := handler.GetEchoContext(c, "HandleChatbackConfirmResourceBorrowed")

	// retrieving item id from payload
	offerItemId := req.Payload.Actions[0].Value
	if offerItemId == nil {
		l.Error("value is required")
		return c.String(http.StatusBadRequest, "value is required")
	}

	// converting item id to item key
	offerItemKey, err := model.ParseOfferItemKey(*offerItemId)
	if err != nil {
		l.Error("could not get offer item id from request", zap.Error(err))
		return c.String(http.StatusBadRequest, err.Error())
	}

	return h.tradingService.ConfirmResourceBorrowed(ctx, offerItemKey)

}

func (h *Handler) HandleChatbackConfirmResourceBorrowedReturned(c echo.Context, req web.InteractionCallback) error {

	ctx, l := handler.GetEchoContext(c, "HandleChatbackConfirmResourceBorrowed")

	// retrieving item id from payload
	offerItemId := req.Payload.Actions[0].Value
	if offerItemId == nil {
		l.Error("value is required")
		return c.String(http.StatusBadRequest, "value is required")
	}

	// converting item id to item key
	offerItemKey, err := model.ParseOfferItemKey(*offerItemId)
	if err != nil {
		l.Error("could not get offer item id from request", zap.Error(err))
		return c.String(http.StatusBadRequest, err.Error())
	}

	return h.tradingService.ConfirmBorrowedResourceReturned(ctx, offerItemKey)

}

func (h *Handler) HandleChatbackOfferAccepted(c echo.Context, req web.InteractionCallback) error {

	ctx, l := handler.GetEchoContext(c, "HandleChatbackOfferAccepted")

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
