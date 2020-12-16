package handler

import (
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/config"
	"github.com/commonpool/backend/pkg/trading"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	chatService    chat.Service
	appConfig      *config.AppConfig
	auth           auth.Authenticator
	tradingService trading.Service
}

func NewHandler(
	chatService chat.Service,
	tradingService trading.Service,
	appConfig *config.AppConfig,
	auth auth.Authenticator) *Handler {
	return &Handler{
		chatService:    chatService,
		tradingService: tradingService,
		appConfig:      appConfig,
		auth:           auth,
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
