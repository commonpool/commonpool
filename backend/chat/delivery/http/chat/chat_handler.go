package chat

import (
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/config"
	"github.com/labstack/echo/v4"
)

type ChatHandler struct {
	chatService chat.Service
	appConfig   *config.AppConfig
	auth        auth.IAuth
}

func NewChatHandler(chatService chat.Service, appConfig *config.AppConfig, auth auth.IAuth) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
		appConfig:   appConfig,
		auth:        auth,
	}
}

func (c *ChatHandler) Register(r *echo.Group) {
	chatGroup := r.Group("/chat", c.auth.Authenticate(true))
	chatGroup.GET("/messages", c.GetMessages)
	chatGroup.GET("/subscriptions", c.GetRecentlyActiveSubscriptions)
	chatGroup.POST("/:id", c.SendMessage)
	chatGroup.POST("/interaction", c.SubmitInteraction)
}
