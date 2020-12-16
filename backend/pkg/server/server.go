package server

import (
	"context"
	"fmt"
	_ "github.com/commonpool/backend/docs"
	"github.com/commonpool/backend/pkg/auth"
	authhandler "github.com/commonpool/backend/pkg/auth/handler"
	"github.com/commonpool/backend/pkg/chat"
	chathandler "github.com/commonpool/backend/pkg/chat/handler"
	chatservice "github.com/commonpool/backend/pkg/chat/service"
	chatstore "github.com/commonpool/backend/pkg/chat/store"
	"github.com/commonpool/backend/pkg/config"
	db2 "github.com/commonpool/backend/pkg/db"
	"github.com/commonpool/backend/pkg/graph"
	grouphandler "github.com/commonpool/backend/pkg/group/handler"
	groupservice "github.com/commonpool/backend/pkg/group/service"
	groupstore "github.com/commonpool/backend/pkg/group/store"
	handler2 "github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/mq"
	"github.com/commonpool/backend/pkg/realtime"
	"github.com/commonpool/backend/pkg/resource"
	resourcehandler "github.com/commonpool/backend/pkg/resource/handler"
	"github.com/commonpool/backend/pkg/resource/service"
	resourcestore "github.com/commonpool/backend/pkg/resource/store"
	"github.com/commonpool/backend/pkg/session"
	"github.com/commonpool/backend/pkg/trading"
	tradingservice "github.com/commonpool/backend/pkg/trading/service"
	tradingstore "github.com/commonpool/backend/pkg/trading/store"
	"github.com/commonpool/backend/pkg/transaction"
	transactionservice "github.com/commonpool/backend/pkg/transaction/service"
	transactionstore "github.com/commonpool/backend/pkg/transaction/store"
	"github.com/commonpool/backend/pkg/user"
	userhandler "github.com/commonpool/backend/pkg/user/handler"
	userservice "github.com/commonpool/backend/pkg/user/service"
	userstore "github.com/commonpool/backend/pkg/user/store"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"time"
)

type Server struct {
	AppConfig          *config.AppConfig
	AmqpClient         mq.Client
	GraphDriver        graph.Driver
	Db                 *gorm.DB
	TransactionStore   transaction.Store
	TransactionService transaction.Service
	UserStore          user.Store
	UserService        user.Service
	ResourceStore      resource.Store
	ResourceService    resource.Service
	ChatStore          chat.Store
	ChatService        chat.Service
	TradingStore       trading.Store
	TradingService     trading.Service
	AuthHandler        *authhandler.AuthHandler
	SessionHandler     *session.Handler
	ChatHandler        *chathandler.Handler
	GroupHandler       *grouphandler.GroupHandler
	ResourceHandler    *resourcehandler.ResourceHandler
	UserHandler        *userhandler.UserHandler
	RealTimeHandler    *realtime.Handler
	Authenticator      auth.Authenticator
	Router             *echo.Echo
}

func NewServer() (*Server, error) {

	ctx := context.Background()

	appConfig, err := config.GetAppConfig(os.LookupEnv, ioutil.ReadFile)
	if err != nil {
		return nil, err
	}

	amqpCli, err := mq.NewRabbitMqClient(ctx, appConfig.AmqpUrl)
	if err != nil {
		return nil, err
	}

	err = graph.InitGraphDatabase(ctx, appConfig)
	if err != nil {
		return nil, err
	}

	driver, err := graph.NewNeo4jDriver(appConfig, appConfig.Neo4jDatabase)
	if err != nil {
		return nil, err
	}

	db := getDb(appConfig)
	db2.AutoMigrate(db)

	transactionStore := transactionstore.NewTransactionStore(db)
	transactionService := transactionservice.NewTransactionService(transactionStore)

	userStore := userstore.NewUserStore(driver)
	userService := userservice.NewUserService(userStore)

	resourceStore := resourcestore.NewResourceStore(driver, transactionService)
	resourceService := service.NewResourceService(resourceStore)

	chatStore := chatstore.NewChatStore(db)
	chatService := chatservice.NewChatService(userStore, amqpCli, chatStore)

	groupStore := groupstore.NewGroupStore(driver)
	groupService := groupservice.NewGroupService(groupStore, amqpCli, chatService, userStore)

	tradingStore := tradingstore.NewTradingStore(driver)
	tradingService := tradingservice.NewTradingService(tradingStore, resourceStore, userStore, chatService, groupService, transactionService)

	r := NewRouter()
	r.HTTPErrorHandler = handler2.HttpErrorHandler
	r.GET("/api/swagger/*", echoSwagger.WrapHandler)

	v1 := r.Group("/api/v1")

	authorization := auth.NewAuth(v1, appConfig, "/api/v1", userStore)

	authHandler := authhandler.NewHandler(authorization)
	authHandler.Register(v1)

	sessionHandler := session.NewHandler(authorization)
	sessionHandler.Register(v1)

	chatHandler := chathandler.NewHandler(chatService, tradingService, appConfig, authorization)
	chatHandler.Register(v1)

	groupHandler := grouphandler.NewHandler(groupService, userService, authorization)
	groupHandler.Register(v1)

	resourceHandler := resourcehandler.NewHandler(resourceService, groupService, userService, authorization)
	resourceHandler.Register(v1)

	userHandler := userhandler.NewHandler(userService, authorization)
	userHandler.Register(v1)

	realtimeHandler := realtime.NewRealtimeHandler(amqpCli, chatService, authorization)
	realtimeHandler.Register(v1)

	return &Server{
		AppConfig:          appConfig,
		AmqpClient:         amqpCli,
		GraphDriver:        driver,
		Db:                 db,
		TransactionStore:   transactionStore,
		TransactionService: transactionService,
		UserStore:          userStore,
		UserService:        userService,
		ResourceStore:      resourceStore,
		ResourceService:    resourceService,
		ChatStore:          chatStore,
		ChatService:        chatService,
		TradingStore:       tradingStore,
		TradingService:     tradingService,
		AuthHandler:        authHandler,
		SessionHandler:     sessionHandler,
		ChatHandler:        chatHandler,
		GroupHandler:       groupHandler,
		ResourceHandler:    resourceHandler,
		UserHandler:        userHandler,
		RealTimeHandler:    realtimeHandler,
		Authenticator:      authorization,
		Router:             r,
	}, nil

}

func (s *Server) Start() error {
	if err := s.Router.Start("0.0.0.0:8585"); err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown() error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.AmqpClient.Shutdown(); err != nil {
		return err
	}

	if err := s.Router.Shutdown(ctx); err != nil {
		return err
	}

	return nil

}

func getDb(appConfig *config.AppConfig) *gorm.DB {
	cs := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", appConfig.DbHost, appConfig.DbUsername, appConfig.DbPassword, appConfig.DbName, appConfig.DbPort)
	db, err := gorm.Open(postgres.Open(cs), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
