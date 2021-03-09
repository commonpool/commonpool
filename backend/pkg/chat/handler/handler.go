package handler

import (
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/chat/service"
	"github.com/commonpool/backend/pkg/config"
	"github.com/commonpool/backend/pkg/trading"
	"github.com/commonpool/backend/pkg/trading/commandhandlers"
	"github.com/commonpool/backend/pkg/trading/queries"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	chatService                service.Service
	appConfig                  *config.AppConfig
	auth                       authenticator.Authenticator
	tradingService             trading.Service
	getOfferKeyForOfferItemKey *queries.GetOfferKeyForOfferItemKey
	confirmBorrowed            *commandhandlers.ConfirmResourceBorrowedHandler
	confirmReturned            *commandhandlers.ConfirmResourceReturnedHandler
	confirmServiceGiven        *commandhandlers.ConfirmServiceGivenHandler
	confirmResourceGiven       *commandhandlers.ConfirmResourceGivenHandler
	declineOffer               *commandhandlers.DeclineOfferHandler
	acceptOffer                *commandhandlers.AcceptOfferHandler
}

func NewHandler(
	chatService service.Service,
	tradingService trading.Service,
	appConfig *config.AppConfig,
	auth authenticator.Authenticator,
	getOfferKeyForOfferItemKey *queries.GetOfferKeyForOfferItemKey,
	confirmBorrowed *commandhandlers.ConfirmResourceBorrowedHandler,
	confirmReturned *commandhandlers.ConfirmResourceReturnedHandler,
	confirmServiceGiven *commandhandlers.ConfirmServiceGivenHandler,
	confirmResourceGiven *commandhandlers.ConfirmResourceGivenHandler,
	declineOffer *commandhandlers.DeclineOfferHandler,
	acceptOffer *commandhandlers.AcceptOfferHandler,
) *Handler {
	return &Handler{
		chatService:                chatService,
		appConfig:                  appConfig,
		auth:                       auth,
		tradingService:             tradingService,
		getOfferKeyForOfferItemKey: getOfferKeyForOfferItemKey,
		confirmBorrowed:            confirmBorrowed,
		confirmReturned:            confirmReturned,
		confirmServiceGiven:        confirmServiceGiven,
		confirmResourceGiven:       confirmResourceGiven,
		declineOffer:               declineOffer,
		acceptOffer:                acceptOffer,
	}
}

func (h *Handler) Register(r *echo.Group) {
	chatGroup := r.Group("/chat", h.auth.Authenticate(true))
	chatGroup.GET("/messages", h.GetMessages)
	chatGroup.GET("/subscriptions", h.GetRecentlyActiveSubscriptions)
	chatGroup.POST("/:id", h.SendMessage)
	chatGroup.POST("/interaction", h.SubmitInteraction)
	chatGroup.POST("/chatback", h.Chatback, h.auth.Authenticate(true))
}
