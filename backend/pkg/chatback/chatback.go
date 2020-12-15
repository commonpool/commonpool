package chatback

import (
	"github.com/commonpool/backend/pkg/auth"
	model2 "github.com/commonpool/backend/pkg/chat/handler/model"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/trading"
	tradingmodel "github.com/commonpool/backend/pkg/trading/model"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	tradingService trading.Service
	authenticator  auth.Authenticator
}

func NewHandler(tradingService trading.Service, authenticator auth.Authenticator) *Handler {
	return &Handler{
		tradingService: tradingService,
		authenticator:  authenticator,
	}
}

func (h *Handler) Register(g *echo.Group) {
	g.POST("/chatback", h.Chatback, h.authenticator.Authenticate(true))
}

func (h *Handler) Chatback(c echo.Context) error {

	// Unmarshal request
	req := model2.InteractionCallback{}
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

func (h *Handler) HandleChatbackConfirmServiceProvided(c echo.Context, req model2.InteractionCallback) error {

	ctx, l := handler.GetEchoContext(c, "HandleChatbackConfirmServiceProvided")

	// retrieving item id from payload
	offerItemId := req.Payload.Actions[0].Value
	if offerItemId == nil {
		l.Error("value is required")
		return c.String(http.StatusBadRequest, "value is required")
	}

	// converting item id to item key
	offerItemKey, err := tradingmodel.ParseOfferItemKey(*offerItemId)
	if err != nil {
		l.Error("could not get offer item id from request", zap.Error(err))
		return c.String(http.StatusBadRequest, err.Error())
	}

	return h.tradingService.ConfirmServiceProvided(ctx, offerItemKey)

}

func (h *Handler) HandleChatbackConfirmResourceTransferred(c echo.Context, req model2.InteractionCallback) error {

	ctx, l := handler.GetEchoContext(c, "HandleChatbackConfirmResourceTransferred")

	// retrieving item id from payload
	offerItemId := req.Payload.Actions[0].Value
	if offerItemId == nil {
		l.Error("value is required")
		return c.String(http.StatusBadRequest, "value is required")
	}

	// converting item id to item key
	offerItemKey, err := tradingmodel.ParseOfferItemKey(*offerItemId)
	if err != nil {
		l.Error("could not get offer item id from request", zap.Error(err))
		return c.String(http.StatusBadRequest, err.Error())
	}

	return h.tradingService.ConfirmResourceTransferred(ctx, offerItemKey)

}

func (h *Handler) HandleChatbackConfirmResourceBorrowed(c echo.Context, req model2.InteractionCallback) error {

	ctx, l := handler.GetEchoContext(c, "HandleChatbackConfirmResourceBorrowed")

	// retrieving item id from payload
	offerItemId := req.Payload.Actions[0].Value
	if offerItemId == nil {
		l.Error("value is required")
		return c.String(http.StatusBadRequest, "value is required")
	}

	// converting item id to item key
	offerItemKey, err := tradingmodel.ParseOfferItemKey(*offerItemId)
	if err != nil {
		l.Error("could not get offer item id from request", zap.Error(err))
		return c.String(http.StatusBadRequest, err.Error())
	}

	return h.tradingService.ConfirmResourceBorrowed(ctx, offerItemKey)

}

func (h *Handler) HandleChatbackConfirmResourceBorrowedReturned(c echo.Context, req model2.InteractionCallback) error {

	ctx, l := handler.GetEchoContext(c, "HandleChatbackConfirmResourceBorrowed")

	// retrieving item id from payload
	offerItemId := req.Payload.Actions[0].Value
	if offerItemId == nil {
		l.Error("value is required")
		return c.String(http.StatusBadRequest, "value is required")
	}

	// converting item id to item key
	offerItemKey, err := tradingmodel.ParseOfferItemKey(*offerItemId)
	if err != nil {
		l.Error("could not get offer item id from request", zap.Error(err))
		return c.String(http.StatusBadRequest, err.Error())
	}

	return h.tradingService.ConfirmBorrowedResourceReturned(ctx, offerItemKey)

}

func (h *Handler) HandleChatbackOfferAccepted(c echo.Context, req model2.InteractionCallback) error {

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

	err = h.tradingService.AcceptOffer(ctx, tradingmodel.NewOfferKey(uid))
	if err != nil {
		l.Error("could not accept offer", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "OK")
}
