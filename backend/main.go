package main

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/config"
	_ "github.com/commonpool/backend/docs"
	"github.com/commonpool/backend/graph"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/handler"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/router"
	"github.com/commonpool/backend/service"
	"github.com/commonpool/backend/store"
	"github.com/commonpool/backend/trading"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
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

	appConfig, err := config.GetAppConfig(os.LookupEnv, ioutil.ReadFile)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	amqpCli, err := amqp.NewRabbitMqClient(ctx, appConfig.AmqpUrl)
	if err != nil {
		log.Fatal(err, "cannot crate amqp client")
	}

	err = graph.InitGraphDatabase(ctx, appConfig)
	if err != nil {
		panic(err)
	}

	driver, err := graph.NewNeo4jDriver(appConfig, appConfig.Neo4jDatabase)
	if err != nil {
		panic(err)
	}

	r := router.NewRouter()

	r.GET("/api/swagger/*", echoSwagger.WrapHandler)

	db := getDb(appConfig)
	store.AutoMigrate(db)

	resourceStore = store.NewResourceStore(driver)
	authStore = store.NewAuthStore(db, driver)
	chatStore := store.NewChatStore(db, authStore, amqpCli)
	tradingStore := store.NewTradingStore(driver)
	groupStore := store.NewGroupStore(driver)

	chatService := service.NewChatService(authStore, groupStore, resourceStore, amqpCli, chatStore)
	groupService := service.NewGroupService(groupStore, amqpCli, chatService, authStore)
	tradingService := service.NewTradingService(tradingStore, resourceStore, authStore, chatService, groupService)

	v1 := r.Group("/api/v1")
	authorization := auth.NewAuth(v1, appConfig, "/api/v1", authStore)

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

	var users []auth.User
	err = db.Model(auth.User{}).Find(&users).Error
	if err != nil {
		panic(err)
	}

	for _, user := range users {
		_, err = chatService.CreateUserExchange(ctx, user.GetUserKey())
		if err != nil {
			panic(err)
		}
	}

	h.Register(v1)

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
		r.Logger.Fatal(err)
	}

	r.Logger.Info("shutting down router")
	if err := r.Shutdown(ctx); err != nil {
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
