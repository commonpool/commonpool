package server

import (
	"context"
	_ "github.com/commonpool/backend/docs"
	"github.com/commonpool/backend/logging"
	"github.com/commonpool/backend/pkg/auth"
	authhandler "github.com/commonpool/backend/pkg/auth/handler"
	chathandler "github.com/commonpool/backend/pkg/chat/handler"
	chatservice "github.com/commonpool/backend/pkg/chat/service"
	chatstore "github.com/commonpool/backend/pkg/chat/store"
	"github.com/commonpool/backend/pkg/chatback"
	"github.com/commonpool/backend/pkg/config"
	db2 "github.com/commonpool/backend/pkg/db"
	"github.com/commonpool/backend/pkg/graph"
	grouphandler "github.com/commonpool/backend/pkg/group/handler"
	groupservice "github.com/commonpool/backend/pkg/group/service"
	groupstore "github.com/commonpool/backend/pkg/group/store"
	handler2 "github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/mq"
	"github.com/commonpool/backend/pkg/realtime"
	resourcehandler "github.com/commonpool/backend/pkg/resource/handler"
	"github.com/commonpool/backend/pkg/resource/service"
	resourcestore "github.com/commonpool/backend/pkg/resource/store"
	"github.com/commonpool/backend/pkg/session"
	tradingservice "github.com/commonpool/backend/pkg/trading/service"
	tradingstore "github.com/commonpool/backend/pkg/trading/store"
	transactionservice "github.com/commonpool/backend/pkg/transaction/service"
	transactionstore "github.com/commonpool/backend/pkg/transaction/store"
	userhandler "github.com/commonpool/backend/pkg/user/handler"
	userservice "github.com/commonpool/backend/pkg/user/service"
	userstore "github.com/commonpool/backend/pkg/user/store"
	"github.com/commonpool/backend/router"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"os"
)

type Server struct {
}

func (s *Server) Start() {

	ctx := context.Background()

	l := logging.WithContext(ctx)

	appConfig, err := config.GetAppConfig(os.LookupEnv, ioutil.ReadFile)
	if err != nil {
		l.Error("could not get app config", zap.Error(err))
		panic(err)
	}

	amqpCli, err := mq.NewRabbitMqClient(ctx, appConfig.AmqpUrl)
	if err != nil {
		log.Fatal(err, "cannot crate amqp client")
	}

	err = graph.InitGraphDatabase(ctx, appConfig)
	if err != nil {
		l.Error("could not initialize graph database", zap.Error(err))
		panic(err)
	}

	driver, err := graph.NewNeo4jDriver(appConfig, appConfig.Neo4jDatabase)
	if err != nil {
		l.Error("could not create neo4j driver", zap.Error(err))
		panic(err)
	}

	db := getDb(appConfig)
	db2.AutoMigrate(db)

	transactionStore := transactionstore.NewTransactionStore(db)
	transactionService := transactionservice.NewTransactionService(transactionStore)
	userStore := userstore.NewUserStore(db, driver)
	resourceStore := resourcestore.NewResourceStore(driver, transactionService)
	groupStore := groupstore.NewGroupStore(driver)
	chatStore := chatstore.NewChatStore(db, userStore, amqpCli)
	tradingStore := tradingstore.NewTradingStore(driver)
	chatService := chatservice.NewChatService(userStore, groupStore, amqpCli, chatStore)
	groupService := groupservice.NewGroupService(groupStore, amqpCli, chatService, userStore)
	userService := userservice.NewUserService(userStore)
	resourceService := service.NewResourceService(resourceStore)
	tradingService := tradingservice.NewTradingService(tradingStore, resourceStore, userStore, chatService, groupService, transactionService)

	r := router.NewRouter()
	r.HTTPErrorHandler = handler2.HttpErrorHandler
	r.GET("/api/swagger/*", echoSwagger.WrapHandler)

	v1 := r.Group("/api/v1")

	authorization := auth.NewAuth(v1, appConfig, "/api/v1", userStore)

	authHandler := authhandler.NewHandler(authorization)
	authHandler.Register(v1)

	sessionHandler := session.NewHandler(authorization)
	sessionHandler.Register(v1)

	chatHandler := chathandler.NewHandler(chatService, appConfig, authorization)
	chatHandler.Register(v1)

	groupHandler := grouphandler.NewHandler(groupService, userService, authorization)
	groupHandler.Register(v1)

	resourceHandler := resourcehandler.NewHandler(resourceService, groupService, userService, authorization)
	resourceHandler.Register(v1)

	chatbackHandler := chatback.NewHandler(tradingService, authorization)
	chatbackHandler.Register(v1)

	userHandler := userhandler.NewHandler(userService, authorization)
	userHandler.Register(v1)

	realtimeHandler := realtime.NewRealtimeHandler(amqpCli, chatService, authorization)
	realtimeHandler.Register(v1)

}
