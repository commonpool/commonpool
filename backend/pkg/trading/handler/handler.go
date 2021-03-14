package handler

import (
	"context"
	"github.com/avast/retry-go"
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	userqueries "github.com/commonpool/backend/pkg/auth/queries"
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	domain2 "github.com/commonpool/backend/pkg/resource/domain"
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
	tradingService                 trading.Service
	groupService                   group.Service
	authorization                  authenticator.Authenticator
	getOfferKeyForOfferItem        *queries.GetOfferKeyForOfferItemKey
	getOffer                       *queries.GetOffer
	getOffers                      *queries.GetOffers
	getOffersWithActions           *queries.GetUserOffersWithActions
	getUsersByKeys                 *userqueries.GetUsersByKeys
	getValueDimensions             *queries.GetValueDimensions
	confirmResourceBorrowedHandler *commandhandlers.ConfirmResourceBorrowedHandler
	confirmResourceReturned        *commandhandlers.ConfirmResourceReturnedHandler
	confirmServiceGiven            *commandhandlers.ConfirmServiceGivenHandler
	confirmResourceGiven           *commandhandlers.ConfirmResourceGivenHandler
	declineOffer                   *commandhandlers.DeclineOfferHandler
	acceptOffer                    *commandhandlers.AcceptOfferHandler
	submitOffer                    *commandhandlers.SubmitOfferHandler
	getGroupReport                 *queries.GetGroupHistory
}

func NewTradingHandler(
	tradingService trading.Service,
	groupService group.Service,
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
	getOffersWithActions *queries.GetUserOffersWithActions,
	getUsersByKeys *userqueries.GetUsersByKeys,
	getValueDimensions *queries.GetValueDimensions,
	getGroupReport *queries.GetGroupHistory,
) *TradingHandler {
	return &TradingHandler{
		tradingService:                 tradingService,
		groupService:                   groupService,
		authorization:                  auth,
		getOfferKeyForOfferItem:        getOfferKeyForOfferItem,
		getOffer:                       getOffer,
		getOffers:                      getOffers,
		getOffersWithActions:           getOffersWithActions,
		getUsersByKeys:                 getUsersByKeys,
		confirmResourceBorrowedHandler: confirmBorrowed,
		confirmResourceReturned:        confirmReturned,
		confirmServiceGiven:            confirmServiceGiven,
		confirmResourceGiven:           confirmResourceGiven,
		declineOffer:                   declineOffer,
		acceptOffer:                    acceptOffer,
		submitOffer:                    submitOffer,
		getValueDimensions:             getValueDimensions,
		getGroupReport:                 getGroupReport,
	}
}

func (h *TradingHandler) Register(e *echo.Group) {
	offers := e.Group("/offers", h.authorization.Authenticate(true))
	offers.GET("/:id", h.HandleGetOffer)
	offers.GET("", h.HandleGetOffers)
	offers.POST("", h.HandleSendOffer)
	offers.POST("/:id/actions/approve", h.HandleAcceptOffer)
	offers.POST("/:id/actions/decline", h.HandleDeclineOffer)
	offers.GET("/target-picker", h.HandleOfferItemTargetPicker)
	offers.POST("/:id/offer-items/:offerItemId/actions/service-given", h.HandleConfirmServiceProvided)
	offers.POST("/:id/offer-items/:offerItemId/actions/resource-given", h.HandleConfirmResourceTransferred)
	offers.POST("/:id/offer-items/:offerItemId/actions/resource-borrowed", h.HandleConfirmResourceBorrowed)
	offers.POST("/:id/offer-items/:offerItemId/actions/resource-returned", h.HandleConfirmBorrowedResourceReturned)
	values := e.Group("/values")
	values.GET("/dimensions", h.HandleGetValueDimensions)
	reports := e.Group("/reports")
	reports.GET("/group-history", h.HandleGetGroupHistory)
}

type SubmitOfferRequest struct {
	Offer SubmitOfferPayload `json:"offer" validate:"required"`
}

type SubmitOfferPayload struct {
	Items    []domain.SubmitOfferItemBase `json:"items" validate:"min=1"`
	GroupKey keys.GroupKey                `json:"groupId" validate:"uuid"`
	Message  string                       `json:"message"`
}

func NewSendOfferPayload(groupKey keys.GroupKeyGetter, items ...domain.SubmitOfferItemBase) SubmitOfferPayload {
	return SubmitOfferPayload{
		Items:    items,
		GroupKey: groupKey.GetGroupKey(),
		Message:  "",
	}
}

func (p SubmitOfferPayload) AsRequest() *SubmitOfferRequest {
	return &SubmitOfferRequest{
		Offer: p,
	}
}

type SendOfferPayloadItem struct {
	Type       domain.OfferItemType `json:"type"`
	To         keys.Target          `json:"to" validate:"required,uuid"`
	From       *keys.Target         `json:"from" validate:"required,uuid"`
	ResourceId *string              `json:"resourceId" validate:"required,uuid"`
	Duration   *string              `json:"duration"`
	Amount     *string              `json:"amount"`
}

type GetOfferResponse struct {
	Offer *groupreadmodels.OfferReadModel `json:"offer"`
}

type GetOffersResponse struct {
	Offers []*groupreadmodels.OfferReadModelWithActions `json:"offers"`
}

func (g GetOfferResponse) GetOfferKey() keys.OfferKey {
	return g.Offer.OfferKey
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
	return c.JSON(http.StatusOK, GetOfferResponse{
		Offer: offer,
	})
}

func (h *TradingHandler) HandleGetOffers(c echo.Context) error {
	ctx, _ := handler.GetEchoContext(c, "HandleGetOffers")
	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	offers, err := h.getOffersWithActions.Get(ctx, loggedInUser.GetUserKey())
	return c.JSON(http.StatusOK, GetOffersResponse{
		Offers: offers,
	})
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

type OfferItemConfirmationRequest struct {
	OfferKey     keys.OfferKey     `param:"id"`
	OfferItemKey keys.OfferItemKey `param:"offerItemId"`
}

func (h *TradingHandler) handleOfferItemAction(c echo.Context, do func(ctx context.Context, offerKey keys.OfferKey, offerItemKey keys.OfferItemKey, loggedInUserKey keys.UserKey) error) error {

	ctx, _ := handler.GetEchoContext(c, "handleOfferItemAction")

	var req OfferItemConfirmationRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}

	if err := do(ctx, req.OfferKey, req.OfferItemKey, loggedInUser.GetUserKey()); err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)

}

func (h *TradingHandler) HandleConfirmBorrowedResourceReturned(c echo.Context) error {
	return h.handleOfferItemAction(c, func(ctx context.Context, offerKey keys.OfferKey, offerItemKey keys.OfferItemKey, loggedInUserKey keys.UserKey) error {
		return h.confirmResourceReturned.Execute(ctx, domain.NewConfirmResourceReturned(ctx, offerKey, offerItemKey, loggedInUserKey))
	})
}

func (h *TradingHandler) HandleConfirmResourceBorrowed(c echo.Context) error {
	return h.handleOfferItemAction(c, func(ctx context.Context, offerKey keys.OfferKey, offerItemKey keys.OfferItemKey, loggedInUserKey keys.UserKey) error {
		return h.confirmResourceBorrowedHandler.Execute(ctx, domain.NewConfirmResourceBorrowed(ctx, offerKey, offerItemKey, loggedInUserKey))
	})
}

func (h *TradingHandler) HandleConfirmResourceTransferred(c echo.Context) error {
	return h.handleOfferItemAction(c, func(ctx context.Context, offerKey keys.OfferKey, offerItemKey keys.OfferItemKey, loggedInUserKey keys.UserKey) error {
		return h.confirmResourceGiven.Execute(ctx, domain.NewConfirmResourceGiven(ctx, offerKey, offerItemKey, loggedInUserKey))
	})
}

func (h *TradingHandler) HandleConfirmServiceProvided(c echo.Context) error {
	return h.handleOfferItemAction(c, func(ctx context.Context, offerKey keys.OfferKey, offerItemKey keys.OfferItemKey, loggedInUserKey keys.UserKey) error {
		return h.confirmServiceGiven.Execute(ctx, domain.NewConfirmServiceGiven(ctx, offerKey, offerItemKey, loggedInUserKey))
	})
}

type GetValueDimensionsResponse struct {
	Dimensions domain2.ValueDimensions `json:"dimensions"`
}

func (h *TradingHandler) HandleGetValueDimensions(c echo.Context) error {
	ctx, _ := handler.GetEchoContext(c, "HandleGetValueDimensions")
	dimensions := h.getValueDimensions.Get(ctx)
	return c.JSON(http.StatusOK, GetValueDimensionsResponse{
		Dimensions: dimensions,
	})
}

func (h *TradingHandler) HandleSendOffer(c echo.Context) error {

	ctx, l := handler.GetEchoContext(c, "HandleSendOffer")

	var err error

	l.Debug("binding request")
	req := SubmitOfferRequest{}
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

	return c.JSON(http.StatusCreated, GetOfferResponse{
		Offer: offer,
	})

}

type GroupReportResponse struct {
	Entries []*groupreadmodels.GroupReportItem `json:"entries"`
}

func (h *TradingHandler) HandleGetGroupHistory(c echo.Context) error {
	ctx, _ := handler.GetEchoContext(c, "HandleSendOffer")

	groupId := c.QueryParam("groupId")
	groupKey, err := keys.ParseGroupKey(groupId)
	if err != nil {
		return err
	}

	histry, err := h.getGroupReport.Get(ctx, groupKey)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &GroupReportResponse{
		Entries: histry,
	})

}
