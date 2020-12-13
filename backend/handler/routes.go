package handler

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(v1 *echo.Group) {

	auth := v1.Group("/auth")
	auth.Any("/login", h.authorization.Login())
	auth.Any("/logout", h.authorization.Logout())

	meta := v1.Group("/meta", h.authorization.Authenticate(false))
	meta.GET("/who-am-i", h.WhoAmI)

	resources := v1.Group("/resources", h.authorization.Authenticate(true))
	resources.GET("", h.SearchResources)
	resources.GET("/:id", h.GetResource)
	resources.POST("", h.CreateResource)
	resources.PUT("/:id", h.UpdateResource)
	resources.POST("/:id/inquire", h.InquireAboutResource)

	users := v1.Group("/users", h.authorization.Authenticate(false))
	users.GET("", h.SearchUsers)
	users.GET("/:id", h.GetUserInfo)
	users.GET("/:id/memberships", h.GetUserMemberships)

	offers := v1.Group("/offers", h.authorization.Authenticate(true))
	offers.GET("/:id", h.HandleGetOffer)
	offers.GET("", h.HandleGetOffers)
	offers.POST("", h.HandleSendOffer)
	offers.POST("/:id/accept", h.HandleAcceptOffer)
	offers.POST("/:id/decline", h.HandleDeclineOffer)
	offers.GET("/target-picker", h.HandleOfferItemTargetPicker)

	tradingHistory := v1.Group("/trading-history", h.authorization.Authenticate(false))
	tradingHistory.POST("", h.GetTradingHistory)

	offerItems := v1.Group("/offer-items", h.authorization.Authenticate(true))
	offerItems.POST("/:id/confirm/service-provided", h.HandleConfirmServiceProvided)
	offerItems.POST("/:id/confirm/resource-transferred", h.HandleConfirmResourceTransferred)
	offerItems.POST("/:id/confirm/resource-borrowed", h.HandleConfirmResourceBorrowed)
	offerItems.POST("/:id/confirm/resource-borrowed-returned", h.HandleConfirmBorrowedResourceReturned)

	my := v1.Group("/my", h.authorization.Authenticate(true))
	my.GET("/memberships", h.GetLoggedInUserMemberships)

	groups := v1.Group("/groups", h.authorization.Authenticate(true))
	groups.POST("", h.CreateGroup)
	groups.GET("/:id", h.GetGroup)
	groups.GET("/:id/memberships", h.GetGroupMemberships)
	groups.GET("/:id/invite-member-picker", h.GetUsersForGroupInvitePicker)
	groups.GET("/:groupId/memberships/:userId", h.GetMembership)

	memberships := v1.Group("/memberships", h.authorization.Authenticate(true))
	memberships.POST("", h.CreateOrAcceptMembership)
	memberships.DELETE("", h.CancelOrDeclineInvitation)

	v1.POST("/chatback", h.Chatback, h.authorization.Authenticate(true))
	v1.GET("/ws", h.Websocket, h.authorization.Authenticate(false))

}
