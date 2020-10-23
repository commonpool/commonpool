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
	users.GET("/:id", h.GetUserInfo)

	meta := v1.Group("/meta", h.authorization.Authenticate(false))
	meta.GET("/who-am-i", h.WhoAmI)

	chat := v1.Group("/chat", h.authorization.Authenticate(true))
	chat.GET("/messages", h.GetMessages)
	chat.GET("/threads", h.GetLatestThreads)
	chat.POST("/:id", h.SendMessage)

	blas := v1.Group("/bla", h.authorization.Authenticate(true))
	blas.GET("", func(context echo.Context) error {
		session := h.authorization.GetAuthUserSession(context)
		if session.Username != "" {
			return context.String(200, "hello, "+session.Username)
		} else {
			return context.String(200, "unauthenticated")
		}

	})
}
