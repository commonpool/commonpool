package main

import (
	"context"
	"fmt"
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
	usermodel "github.com/commonpool/backend/pkg/user/model"
	userservice "github.com/commonpool/backend/pkg/user/service"
	userstore "github.com/commonpool/backend/pkg/user/store"
	"github.com/commonpool/backend/router"
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

	var users []usermodel.User
	err = db.Model(usermodel.User{}).Find(&users).Error
	if err != nil {
		l.Error("could not find users", zap.Error(err))
		panic(err)
	}

	for _, u := range users {
		_, err = chatService.CreateUserExchange(ctx, u.GetUserKey())
		if err != nil {
			l.Error("could not create user exchange for user", zap.Object("user", u.GetUserKey()), zap.Error(err))
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
