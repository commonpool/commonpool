package main

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/config"
	_ "github.com/commonpool/backend/docs"
	"github.com/commonpool/backend/graph"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/handler"
	"github.com/commonpool/backend/logging"
	"github.com/commonpool/backend/pkg/chat"
	chathandler "github.com/commonpool/backend/pkg/chat/handler"
	chatservice "github.com/commonpool/backend/pkg/chat/service"
	chatstore "github.com/commonpool/backend/pkg/chat/store"
	handler2 "github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/router"
	"github.com/commonpool/backend/service"
	"github.com/commonpool/backend/store"
	"github.com/commonpool/backend/trading"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"
)

var (
	d             *gorm.DB
	resourceStore resource.Store
	authStore     auth.Store
	cs            chat.Store
	ts            trading.Store
	gs            group.Store
	e             *echo.Echo
)

/*
{"neo4j":"FOLLOWER","system":"FOLLOWER","bla":"LEADER"}
*/

// @title commonpool api
// @version 1.0
// @description resources api
// @termsOfService http://swagger.io/terms
// @host 127.0.0.1:8585
// @basePath /api/v1
func main() {

	ctx := context.Background()
	l := logging.WithContext(ctx)

	appConfig, err := config.GetAppConfig(os.LookupEnv, ioutil.ReadFile)
	if err != nil {
		l.Error("could not get app config", zap.Error(err))
		panic(err)
	}

	amqpCli, err := amqp.NewRabbitMqClient(ctx, appConfig.AmqpUrl)
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

	r := router.NewRouter()

	r.GET("/api/swagger/*", echoSwagger.WrapHandler)

	db := getDb(appConfig)
	store.AutoMigrate(db)

	transactionStore := store.NewTransactionStore(db)
	transactionService := service.NewTransactionService(transactionStore)
	resourceStore = store.NewResourceStore(driver, transactionService)
	authStore = store.NewAuthStore(db, driver)
	chatStore := chatstore.NewChatStore(db, authStore, amqpCli)
	tradingStore := store.NewTradingStore(driver)
	groupStore := store.NewGroupStore(driver)

	chatService := chatservice.NewChatService(authStore, groupStore, resourceStore, amqpCli, chatStore)
	groupService := service.NewGroupService(groupStore, amqpCli, chatService, authStore)
	tradingService := service.NewTradingService(tradingStore, resourceStore, authStore, chatService, groupService, transactionService)

	v1 := r.Group("/api/v1")
	authorization := auth.NewAuth(v1, appConfig, "/api/v1", authStore)

	r.HTTPErrorHandler = handler2.HttpErrorHandler

	chatHandler := chathandler.NewChatHandler(chatService, appConfig, authorization)
	chatHandler.Register(v1)

	h := handler.NewHandler(
		resourceStore,
		authStore,
		chatStore,
		tradingStore,
		authorization,
		amqpCli,
		*appConfig,
		chatService,
		tradingService,
		groupService)

	h.Register(v1)

	var users []auth.User
	err = db.Model(auth.User{}).Find(&users).Error
	if err != nil {
		l.Error("could not find users", zap.Error(err))
		panic(err)
	}

	for _, user := range users {
		_, err = chatService.CreateUserExchange(ctx, user.GetUserKey())
		if err != nil {
			l.Error("could not create user exchange for user", zap.Object("user", user.GetUserKey()), zap.Error(err))
			panic(err)
		}
	}

	// Start server
	go func() {
		if err := r.Start("0.0.0.0:8585"); err != nil {
			r.Logger.Error(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r.Logger.Info("shutting down amqp client")
	if err := amqpCli.Shutdown(); err != nil {
		l.Error("could nots shutdown amqp client", zap.Error(err))
		r.Logger.Fatal(err)
	}

	r.Logger.Info("shutting down router")
	if err := r.Shutdown(ctx); err != nil {
		l.Error("could not shut down router", zap.Error(err))
		r.Logger.Fatal(err)
	}

}

func getDb(appConfig *config.AppConfig) *gorm.DB {
	cs := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", appConfig.DbHost, appConfig.DbUsername, appConfig.DbPassword, appConfig.DbName, appConfig.DbPort)
	db, err := gorm.Open(postgres.Open(cs), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
