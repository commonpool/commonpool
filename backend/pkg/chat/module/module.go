package module

import (
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/auth/store"
	chathandler "github.com/commonpool/backend/pkg/chat/handler"
	service2 "github.com/commonpool/backend/pkg/chat/service"
	store2 "github.com/commonpool/backend/pkg/chat/store"
	"github.com/commonpool/backend/pkg/config"
	"github.com/commonpool/backend/pkg/mq"
	"github.com/commonpool/backend/pkg/trading/service"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Module struct {
	Store   store2.Store
	Service service2.Service
	Handler *chathandler.Handler
}

func NewModule(appConfig *config.AppConfig, authenticator authenticator.Authenticator, amqpClient mq.Client, userStore store.Store, db *gorm.DB, tradingService service.TradingService) *Module {
	chatStore := store2.NewChatStore(db)
	chatService := service2.NewChatService(userStore, amqpClient, chatStore)
	chatHandler := chathandler.NewHandler(chatService, tradingService, appConfig, authenticator)
	return &Module{
		Store:   chatStore,
		Service: chatService,
		Handler: chatHandler,
	}
}

func (m *Module) Register(g *echo.Group) {
	m.Handler.Register(g)
}
