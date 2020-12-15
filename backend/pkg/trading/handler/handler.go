package handler

import (
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/trading"
	"github.com/commonpool/backend/pkg/user"
	"github.com/labstack/echo/v4"
)

type TradingHandler struct {
	tradingService trading.Service
	groupService   group.Service
	userService    user.Service
	authorization  auth.Authenticator
}

func NewTradingHandler(
	tradingService trading.Service,
	groupService group.Service,
	userService user.Service,
	auth auth.Authenticator) *TradingHandler {
	return &TradingHandler{
		tradingService: tradingService,
		groupService:   groupService,
		userService:    userService,
		authorization:  auth,
	}
}

func (h *TradingHandler) Register(e *echo.Echo) {

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
