package handler

import (
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/auth/service"
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/trading"
	"github.com/commonpool/backend/pkg/trading/domain"
	"github.com/commonpool/backend/pkg/trading/queries"
	"github.com/labstack/echo/v4"
)

type TradingHandler struct {
	tradingService          trading.Service
	groupService            group.Service
	userService             service.Service
	authorization           authenticator.Authenticator
	offerRepo               domain.OfferRepository
	getOfferKeyForOfferItem *queries.GetOfferKeyForOfferItemKey
	getOfferItem            *queries.GetOfferItem
	getOffer                *queries.GetOffer
	getPermissions          *queries.GetOfferPermissions
}

func NewTradingHandler(
	tradingService trading.Service,
	groupService group.Service,
	userService service.Service,
	auth authenticator.Authenticator,
	offerRepo domain.OfferRepository,
	getOfferKeyForOfferItem *queries.GetOfferKeyForOfferItemKey,
	getOfferItem *queries.GetOfferItem,
	getOffer *queries.GetOffer,
	getPermissions *queries.GetOfferPermissions,
) *TradingHandler {
	return &TradingHandler{
		tradingService:          tradingService,
		groupService:            groupService,
		userService:             userService,
		authorization:           auth,
		offerRepo:               offerRepo,
		getOfferKeyForOfferItem: getOfferKeyForOfferItem,
		getOfferItem:            getOfferItem,
		getOffer:                getOffer,
		getPermissions:          getPermissions,
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
