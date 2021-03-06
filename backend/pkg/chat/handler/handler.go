package handler

import (
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/chat/service"
	"github.com/commonpool/backend/pkg/config"
	"github.com/commonpool/backend/pkg/trading"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	chatService    service.Service
	appConfig      *config.AppConfig
	auth           authenticator.Authenticator
	tradingService trading.Service
}

func NewHandler(
	chatService service.Service,
	tradingService trading.Service,
	appConfig *config.AppConfig,
	auth authenticator.Authenticator) *Handler {
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
