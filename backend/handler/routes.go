package handler

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(v1 *echo.Group) {

	auth := v1.Group("/auth")
	auth.Any("/login", h.authorization.Login())
	auth.Any("/logout", h.authorization.Logout())

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

	meta := v1.Group("/meta", h.authorization.Authenticate(false))
	meta.GET("/who-am-i", h.WhoAmI)

	chat := v1.Group("/chat", h.authorization.Authenticate(true))
	chat.GET("/messages", h.GetMessages)
	chat.GET("/subscriptions", h.GetRecentlyActiveSubscriptions)
	chat.POST("/:id", h.SendMessage)
	chat.POST("/interaction", h.SubmitInteraction)

	offers := v1.Group("/offers", h.authorization.Authenticate(true))
	offers.POST("", h.HandleSendOffer)
	offers.GET("/:id", h.GetOffer)
	offers.GET("", h.GetOffers)
	offers.POST("/:id/accept", h.HandleAcceptOffer)
	offers.POST("/:id/decline", h.DeclineOffer)

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
	v1.GET("/ws", h.Websocket, h.authorization.Authenticate(true))

}
