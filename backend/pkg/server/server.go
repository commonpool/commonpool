package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/bsm/redislock"
	_ "github.com/commonpool/backend/docs"
	"github.com/commonpool/backend/pkg/auth"
	authhandler "github.com/commonpool/backend/pkg/auth/handler"
	"github.com/commonpool/backend/pkg/chat"
	chathandler "github.com/commonpool/backend/pkg/chat/handler"
	chatservice "github.com/commonpool/backend/pkg/chat/service"
	chatstore "github.com/commonpool/backend/pkg/chat/store"
	"github.com/commonpool/backend/pkg/clusterlock"
	"github.com/commonpool/backend/pkg/config"
	db2 "github.com/commonpool/backend/pkg/db"
	"github.com/commonpool/backend/pkg/eventbus"
	postgres2 "github.com/commonpool/backend/pkg/eventstore/postgres"
	"github.com/commonpool/backend/pkg/eventstore/publish"
	"github.com/commonpool/backend/pkg/graph"
	grouphandler "github.com/commonpool/backend/pkg/group/handler"
	groupservice "github.com/commonpool/backend/pkg/group/service"
	groupstore "github.com/commonpool/backend/pkg/group/store"
	handler2 "github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/mq"
	nukehandler "github.com/commonpool/backend/pkg/nuke/handler"
	"github.com/commonpool/backend/pkg/realtime"
	"github.com/commonpool/backend/pkg/resource"
	resourcehandler "github.com/commonpool/backend/pkg/resource/handler"
	"github.com/commonpool/backend/pkg/resource/service"
	resourcestore "github.com/commonpool/backend/pkg/resource/store"
	"github.com/commonpool/backend/pkg/session"
	"github.com/commonpool/backend/pkg/trading"
	tradinghandler "github.com/commonpool/backend/pkg/trading/handler"
	"github.com/commonpool/backend/pkg/trading/listeners"
	tradingservice "github.com/commonpool/backend/pkg/trading/service"
	tradingstore "github.com/commonpool/backend/pkg/trading/store"
	"github.com/commonpool/backend/pkg/transaction"
	transactionservice "github.com/commonpool/backend/pkg/transaction/service"
	transactionstore "github.com/commonpool/backend/pkg/transaction/store"
	"github.com/commonpool/backend/pkg/user"
	userhandler "github.com/commonpool/backend/pkg/user/handler"
	userservice "github.com/commonpool/backend/pkg/user/service"
	userstore "github.com/commonpool/backend/pkg/user/store"
	"github.com/go-redis/redis/v8"
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
	GroupHandler       *grouphandler.Handler
	ResourceHandler    *resourcehandler.ResourceHandler
	UserHandler        *userhandler.UserHandler
	RealTimeHandler    *realtime.Handler
	Authenticator      auth.Authenticator
	Router             *echo.Echo
	NukeHandler        *nukehandler.Handler
	TradingHandler     *tradinghandler.TradingHandler
	RedisClient        *redis.Client
	ClusterLocker      clusterlock.Locker
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func NewServer() (*Server, error) {

	ctx := context.Background()

	appConfig, err := config.GetAppConfig(os.LookupEnv, ioutil.ReadFile)
	if err != nil {
		return nil, err
	}

	var redisTlsConfig *tls.Config = nil
	if appConfig.RedisTlsEnabled {
		redisTlsConfig = &tls.Config{
			InsecureSkipVerify: appConfig.RedisTlsSkipVerify,
		}
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:      appConfig.RedisHost + ":" + appConfig.RedisPort,
		Password:  appConfig.RedisPassword,
		DB:        0,
		TLSConfig: redisTlsConfig,
	})

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	clusterLocker := clusterlock.NewRedis(redislock.New(redisClient))

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

	eventPublisher := eventbus.NewAmqpPublisher(amqpCli)
	eventStore := publish.NewPublishEventStore(postgres2.NewPostgresEventStore(db), eventPublisher)

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

	offerRepository := tradingstore.NewEventSourcedOfferRepository(eventStore)

	tradingStore := tradingstore.NewTradingStore(driver)
	tradingService := tradingservice.NewTradingService(
		tradingStore,
		resourceStore,
		userStore,
		chatService,
		groupService,
		transactionService,
		offerRepository)

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

	tradingHandler := tradinghandler.NewTradingHandler(tradingService, groupService, userService, authorization)
	tradingHandler.Register(v1)

	nukeHandler := nukehandler.NewHandler(db, amqpCli, driver)
	nukeHandler.Register(v1)

	var catchUpListenerFactory eventbus.CatchUpListenerFactory = func(key string, lockTTL time.Duration) *eventbus.CatchUpListener {
		return eventbus.NewCatchUpListener(
			eventStore,
			func() time.Time { return time.Time{} },
			amqpCli,
			eventbus.NewRedisDeduplicator(100, redisClient, key),
			clusterLocker,
			lockTTL,
			&clusterlock.Options{
				RetryStrategy: clusterlock.EverySecond,
			},
		)
	}

	handler := listeners.NewTransactionHistoryHandler(db, catchUpListenerFactory)
	go func() {
		handler.Start(ctx)
	}()

	offerRm := listeners.NewOfferReadModelHandler(db, catchUpListenerFactory)
	go func() {
		offerRm.Start(ctx)
	}()

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
		NukeHandler:        nukeHandler,
		Authenticator:      authorization,
		TradingHandler:     tradingHandler,
		Router:             r,
		ClusterLocker:      clusterLocker,
		RedisClient:        redisClient,
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
