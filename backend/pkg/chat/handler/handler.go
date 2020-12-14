package handler

import (
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/config"
	"github.com/labstack/echo/v4"
)

type ChatHandler struct {
	chatService chat.Service
	appConfig   *config.AppConfig
	auth        auth.Authenticator
}

func NewChatHandler(
	chatService chat.Service,
	appConfig *config.AppConfig,
	auth auth.Authenticator) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
		appConfig:   appConfig,
		auth:        auth,
	}
}

func (chatHandler *ChatHandler) Register(r *echo.Group) {
	chatGroup := r.Group("/chat", chatHandler.auth.Authenticate(true))
	chatGroup.GET("/messages", chatHandler.GetMessages)
	chatGroup.GET("/subscriptions", chatHandler.GetRecentlyActiveSubscriptions)
	chatGroup.POST("/:id", chatHandler.SendMessage)
	chatGroup.POST("/interaction", chatHandler.SubmitInteraction)
}
