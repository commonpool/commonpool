package handler

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(v1 *echo.Group) {
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
	chat.GET("/threads", h.GetLatestThreads)
	chat.POST("/:id", h.SendMessage)

	offers := v1.Group("/offers", h.authorization.Authenticate(true))
	offers.POST("", h.SendOffer)
	offers.GET("/:id", h.GetOffer)
	offers.GET("", h.GetOffers)
	offers.POST("/:id/accept", h.AcceptOffer)
	offers.POST("/:id/decline", h.DeclineOffer)

	my := v1.Group("/my", h.authorization.Authenticate(true))
	my.GET("/memberships", h.GetLoggedInUserMemberships)

	groups := v1.Group("/groups", h.authorization.Authenticate(true))
	groups.POST("", h.CreateGroup)
	groups.GET("/:id", h.GetGroup)
	groups.GET("/:id/memberships", h.GetGroupMemberships)
	groups.GET("/:id/invite-member-picker", h.GetUsersForGroupInvitePicker)
	groups.POST("/:id/invite", h.InviteUser)
	groups.GET("/:groupId/memberships/:userId", h.GetMembership)
	groups.POST("/:groupId/memberships/:userId/accept", h.AcceptInvitation)
	groups.POST("/:groupId/memberships/:userId/decline", h.DeclineInvitation)
	groups.DELETE("/:groupId/memberships/:userId", h.LeaveGroup)

}
