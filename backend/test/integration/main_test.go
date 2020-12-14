package integration

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/config"
	"github.com/commonpool/backend/handler"
	"github.com/commonpool/backend/mock"
	chatservice "github.com/commonpool/backend/pkg/chat/service"
	chatstore "github.com/commonpool/backend/pkg/chat/store"
	"github.com/commonpool/backend/pkg/db"
	graph2 "github.com/commonpool/backend/pkg/graph"
	groupservice "github.com/commonpool/backend/pkg/group/service"
	groupstore "github.com/commonpool/backend/pkg/group/store"
	resourcestore "github.com/commonpool/backend/pkg/resource/store"
	tradingservice "github.com/commonpool/backend/pkg/trading/service"
	tradingstore "github.com/commonpool/backend/pkg/trading/store"
	transactionservice "github.com/commonpool/backend/pkg/transaction/service"
	transactionstore "github.com/commonpool/backend/pkg/transaction/store"
	userstore "github.com/commonpool/backend/pkg/user/store"
	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	_ "net/http/pprof"
	"os"
	"sync"
	"testing"
)

var a *handler.Handler
var Db *gorm.DB
var AmqpClient amqp.Client
var ResourceStore resourcestore.ResourceStore
var AuthStore userstore.UserStore
var ChatStore chatstore.ChatStore
var TradingStore tradingstore.TradingStore
var GroupStore groupstore.GroupStore
var ChatService chatservice.ChatService
var TradingService tradingservice.TradingService
var GroupService groupservice.GroupService
var Authorizer *mock.AuthenticatorMock
var Driver *graph2.Neo4jGraphDriver
var TransactionStore *transactionstore.TransactionStore
var TransactionService *transactionservice.TransactionService

func TestMain(m *testing.M) {

	println("running main")

	appConfig, err := config.GetAppConfig(os.LookupEnv, ioutil.ReadFile)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	AmqpClient, err = amqp.NewRabbitMqClient(ctx, appConfig.AmqpUrl)
	if err != nil {
		panic(err)
	}
	Authorizer = &mock.AuthenticatorMock{}

	err = graph2.InitGraphDatabase(nil, appConfig)
	if err != nil {
		panic(err)
	}

	Driver, err = graph2.NewNeo4jDriver(appConfig, appConfig.Neo4jDatabase)
	if err != nil {
		panic(err)
	}

	Db = getDb(appConfig)

	TransactionStore = transactionstore.NewTransactionStore(Db)
	TransactionService = transactionservice.NewTransactionService(TransactionStore)
	ResourceStore = *resourcestore.NewResourceStore(Driver, TransactionService)
	AuthStore = *userstore.NewAuthStore(Db, Driver)
	ChatStore = *chatstore.NewChatStore(Db, &AuthStore, AmqpClient)
	TradingStore = *tradingstore.NewTradingStore(Driver)
	GroupStore = *groupstore.NewGroupStore(Driver)
	ChatService = *chatservice.NewChatService(&AuthStore, &GroupStore, &ResourceStore, AmqpClient, &ChatStore)
	GroupService = *groupservice.NewGroupService(&GroupStore, AmqpClient, ChatService, &AuthStore)
	TradingService = *tradingservice.NewTradingService(TradingStore, &ResourceStore, &AuthStore, ChatService, GroupService, TransactionService)

	db.AutoMigrate(Db)

	a = handler.NewHandler(
		&ResourceStore,
		&AuthStore,
		&ChatStore,
		TradingStore,
		Authorizer,
		AmqpClient,
		*appConfig,
		ChatService,
		TradingService,
		GroupService)

	cleanDb()

	session, err := Driver.GetSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	result, err := session.Run(`MATCH (u:User) detach delete u`, map[string]interface{}{})
	if err != nil {
		panic(err)
	}
	if result.Err() != nil {
		panic(result.Err())
	}

	os.Exit(m.Run())

}

var userIncrementer = 0
var userIncrementerMu sync.Mutex

func NewUser() *auth.UserSession {
	userIncrementerMu.Lock()
	defer func() {
		userIncrementerMu.Unlock()
	}()
	userIncrementer++
	var userId = uuid.NewV4().String()
	userEmail := fmt.Sprintf("user%d@email.com", userIncrementer)
	userName := fmt.Sprintf("user%d", userIncrementer)
	return &auth.UserSession{
		Username:        userName,
		Subject:         userId,
		Email:           userEmail,
		IsAuthenticated: true,
	}
}

var createUserLock sync.Mutex

func cleanDb() {

	session, err := Driver.GetSession()
	if err != nil {
		panic(err)
	}

	_, err = session.Run(`MATCH (n) DETACH DELETE n`, map[string]interface{}{})
	if err != nil {
		panic(err)
	}

	Db.Delete(chatstore.Channel{}, "1 = 1")
	Db.Delete(chatstore.ChannelSubscription{}, "1 = 1")
	Db.Delete(chatstore.Message{}, "1 = 1")
}

func getDb(appConfig *config.AppConfig) *gorm.DB {
	cs := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", appConfig.DbHost, appConfig.DbUsername, appConfig.DbPassword, appConfig.DbName, appConfig.DbPort)
	database, err := gorm.Open(postgres.Open(cs), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return database
}
