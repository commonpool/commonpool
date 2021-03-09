package handler

import (
	"github.com/avast/retry-go"
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/auth/service"
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading"
	"github.com/commonpool/backend/pkg/trading/commandhandlers"
	"github.com/commonpool/backend/pkg/trading/domain"
	"github.com/commonpool/backend/pkg/trading/queries"
	groupreadmodels "github.com/commonpool/backend/pkg/trading/readmodels"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type TradingHandler struct {
	tradingService          trading.Service
	groupService            group.Service
	userService             service.Service
	authorization           authenticator.Authenticator
	getOfferKeyForOfferItem *queries.GetOfferKeyForOfferItemKey
	getOffer                *queries.GetOffer
	getOffers               *queries.GetOffers
	confirmResourceBorrowed *commandhandlers.ConfirmResourceBorrowedHandler
	confirmResourceReturned *commandhandlers.ConfirmResourceReturnedHandler
	confirmServiceGiven     *commandhandlers.ConfirmServiceGivenHandler
	confirmResourceGiven    *commandhandlers.ConfirmResourceGivenHandler
	declineOffer            *commandhandlers.DeclineOfferHandler
	acceptOffer             *commandhandlers.AcceptOfferHandler
	submitOffer             *commandhandlers.SubmitOfferHandler
}

func NewTradingHandler(
	tradingService trading.Service,
	groupService group.Service,
	userService service.Service,
	auth authenticator.Authenticator,
	getOffer *queries.GetOffer,
	getOffers *queries.GetOffers,
	getOfferKeyForOfferItem *queries.GetOfferKeyForOfferItemKey,
	confirmBorrowed *commandhandlers.ConfirmResourceBorrowedHandler,
	confirmReturned *commandhandlers.ConfirmResourceReturnedHandler,
	confirmServiceGiven *commandhandlers.ConfirmServiceGivenHandler,
	confirmResourceGiven *commandhandlers.ConfirmResourceGivenHandler,
	declineOffer *commandhandlers.DeclineOfferHandler,
	acceptOffer *commandhandlers.AcceptOfferHandler,
	submitOffer *commandhandlers.SubmitOfferHandler,
) *TradingHandler {
	return &TradingHandler{
		tradingService:          tradingService,
		groupService:            groupService,
		userService:             userService,
		authorization:           auth,
		getOfferKeyForOfferItem: getOfferKeyForOfferItem,
		getOffer:                getOffer,
		getOffers:               getOffers,
		confirmResourceBorrowed: confirmBorrowed,
		confirmResourceReturned: confirmReturned,
		confirmServiceGiven:     confirmServiceGiven,
		confirmResourceGiven:    confirmResourceGiven,
		declineOffer:            declineOffer,
		acceptOffer:             acceptOffer,
		submitOffer:             submitOffer,
	}
}

func (h *TradingHandler) Register(e *echo.Group) {
	offers := e.Group("/offers", h.authorization.Authenticate(true))
	offers.GET("/:id", h.HandleGetOffer)
	offers.GET("", h.HandleGetOffers)
	offers.POST("", h.HandleSendOffer)
	offers.POST("/:id/accept", h.HandleAcceptOffer)
	offers.POST("/:id/decline", h.HandleDeclineOffer)
	offers.GET("/target-picker", h.HandleOfferItemTargetPicker)
	offerItems := e.Group("/offer-items", h.authorization.Authenticate(true))
	offerItems.POST("/:id/confirm/service-provided", h.HandleConfirmServiceProvided)
	offerItems.POST("/:id/confirm/resource-transferred", h.HandleConfirmResourceTransferred)
	offerItems.POST("/:id/confirm/resource-borrowed", h.HandleConfirmResourceBorrowed)
	offerItems.POST("/:id/confirm/resource-borrowed-returned", h.HandleConfirmBorrowedResourceReturned)
}

type SendOfferRequest struct {
	Offer SendOfferPayload `json:"offer" validate:"required"`
}

type SendOfferPayload struct {
	Items    []domain.SubmitOfferItemBase `json:"items" validate:"min=1"`
	GroupKey keys.GroupKey                `json:"groupId" validate:"uuid"`
	Message  string                       `json:"message"`
}

type SendOfferPayloadItem struct {
	Type       domain.OfferItemType `json:"type"`
	To         keys.Target          `json:"to" validate:"required,uuid"`
	From       *keys.Target         `json:"from" validate:"required,uuid"`
	ResourceId *string              `json:"resourceId" validate:"required,uuid"`
	Duration   *string              `json:"duration"`
	Amount     *string              `json:"amount"`
}

type OfferResponse struct {
	Offer *groupreadmodels.OfferReadModel
}

func (h *TradingHandler) HandleGetOffer(c echo.Context) error {
	ctx := handler.GetContext(c)
	offerKey, err := keys.ParseOfferKey(c.Param("id"))
	if err != nil {
		return err
	}
	offer, err := h.getOffer.Get(ctx, offerKey)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, offer)
}

func (h *TradingHandler) HandleGetOffers(c echo.Context) error {
	ctx, _ := handler.GetEchoContext(c, "HandleGetOffers")
	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	offers, err := h.getOffers.Get(ctx, loggedInUser.GetUserKey())
	return c.JSON(http.StatusOK, offers)
}

func (h *TradingHandler) HandleAcceptOffer(c echo.Context) error {
	ctx, _ := handler.GetEchoContext(c, "HandleAcceptOffer")
	offerKey, err := keys.ParseOfferKey(c.Param("id"))
	if err != nil {
		return err
	}
	if err := h.acceptOffer.Execute(ctx, domain.NewAcceptOffer(ctx, offerKey)); err != nil {
		return err
	}
	return c.NoContent(http.StatusAccepted)
}

func (h *TradingHandler) HandleDeclineOffer(c echo.Context) error {
	ctx, _ := handler.GetEchoContext(c, "HandleDeclineOffer")
	offerKey, err := keys.ParseOfferKey(c.Param("id"))
	if err != nil {
		return err
	}
	if err := h.declineOffer.Execute(ctx, domain.NewDeclineOffer(ctx, offerKey)); err != nil {
		return err
	}
	return c.NoContent(http.StatusAccepted)
}

func (h *TradingHandler) HandleConfirmBorrowedResourceReturned(c echo.Context) error {
	ctx, _ := handler.GetEchoContext(c, "HandleConfirmBorrowedResourceReturned")
	offerItemKey, err := keys.ParseOfferItemKey(c.Param("id"))
	if err != nil {
		return err
	}
	offerKey, err := h.getOfferKeyForOfferItem.Get(ctx, offerItemKey)
	if err != nil {
		return err
	}
	if err := h.confirmResourceReturned.Execute(ctx, domain.NewConfirmResourceReturned(ctx, offerKey, offerItemKey)); err != nil {
		return err
	}
	return c.NoContent(http.StatusAccepted)
}

func (h *TradingHandler) HandleConfirmResourceBorrowed(c echo.Context) error {
	ctx, _ := handler.GetEchoContext(c, "HandleConfirmResourceBorrowed")
	offerItemKey, err := keys.ParseOfferItemKey(c.Param("id"))
	if err != nil {
		return err
	}
	offerKey, err := h.getOfferKeyForOfferItem.Get(ctx, offerItemKey)
	if err != nil {
		return err
	}
	if err := h.confirmResourceBorrowed.Execute(ctx, domain.NewConfirmResourceBorrowed(ctx, offerKey, offerItemKey)); err != nil {
		return err
	}
	return c.NoContent(http.StatusAccepted)
}

func (h *TradingHandler) HandleConfirmResourceTransferred(c echo.Context) error {
	ctx, _ := handler.GetEchoContext(c, "HandleConfirmResourceTransferred")
	offerItemKey, err := keys.ParseOfferItemKey(c.Param("id"))
	if err != nil {
		return err
	}
	offerKey, err := h.getOfferKeyForOfferItem.Get(ctx, offerItemKey)
	if err != nil {
		return err
	}
	if err := h.confirmResourceGiven.Execute(ctx, domain.NewConfirmResourceGiven(ctx, offerKey, offerItemKey)); err != nil {
		return err
	}
	return c.NoContent(http.StatusAccepted)
}

func (h *TradingHandler) HandleConfirmServiceProvided(c echo.Context) error {
	ctx, _ := handler.GetEchoContext(c, "HandleConfirmServiceProvided")
	offerItemKey, err := keys.ParseOfferItemKey(c.Param("id"))
	if err != nil {
		return err
	}
	offerKey, err := h.getOfferKeyForOfferItem.Get(ctx, offerItemKey)
	if err != nil {
		return err
	}
	if err := h.confirmServiceGiven.Execute(ctx, domain.NewConfirmServiceGiven(ctx, offerKey, offerItemKey)); err != nil {
		return err
	}
	return c.NoContent(http.StatusAccepted)
}

func (h *TradingHandler) HandleSendOffer(c echo.Context) error {

	ctx, l := handler.GetEchoContext(c, "HandleSendOffer")

	var err error

	l.Debug("binding request")
	req := SendOfferRequest{}
	if err = c.Bind(&req); err != nil {
		l.Error("could not bind request", zap.Error(err))
		return err
	}

	l.Debug("validating request")
	if err = c.Validate(req); err != nil {
		l.Error("error validating request", zap.Error(err))
		return err
	}

	l.Debug("posting offer")
	offerKey := keys.GenerateOfferKey()

	var mappedItems []domain.SubmitOfferItem
	for _, item := range req.Offer.Items {
		mappedItems = append(mappedItems, domain.SubmitOfferItem{
			SubmitOfferItemBase: item,
			OfferItemKey:        keys.GenerateOfferItemKey(),
		})
	}

	err = h.submitOffer.Execute(ctx, domain.NewPostOffer(ctx, offerKey, req.Offer.GroupKey, mappedItems))
	if err != nil {
		return err
	}

	l.Debug("getting offer readmodel")
	var offer *groupreadmodels.OfferReadModel
	err = retry.Do(func() error {
		offer, err = h.getOffer.Get(ctx, offerKey)
		return err
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, OfferResponse{
		Offer: offer,
	})

}
