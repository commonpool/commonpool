package integration

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/mock"
	authhandler "github.com/commonpool/backend/pkg/auth/handler"
	"github.com/commonpool/backend/pkg/auth/models"
	service2 "github.com/commonpool/backend/pkg/auth/service"
	"github.com/commonpool/backend/pkg/auth/store"
	chathandler "github.com/commonpool/backend/pkg/chat/handler"
	chatservice "github.com/commonpool/backend/pkg/chat/service"
	chatstore "github.com/commonpool/backend/pkg/chat/store"
	"github.com/commonpool/backend/pkg/config"
	"github.com/commonpool/backend/pkg/db"
	graph2 "github.com/commonpool/backend/pkg/graph"
	grouphandler "github.com/commonpool/backend/pkg/group/handler"
	groupservice "github.com/commonpool/backend/pkg/group/service"
	groupstore "github.com/commonpool/backend/pkg/group/store"
	"github.com/commonpool/backend/pkg/mq"
	resourcehandler "github.com/commonpool/backend/pkg/resource/handler"
	"github.com/commonpool/backend/pkg/resource/service"
	resourcestore "github.com/commonpool/backend/pkg/resource/store"
	"github.com/commonpool/backend/pkg/session"
	tradinghandler "github.com/commonpool/backend/pkg/trading/handler"
	tradingservice "github.com/commonpool/backend/pkg/trading/service"
	tradingstore "github.com/commonpool/backend/pkg/trading/store"
	transactionservice "github.com/commonpool/backend/pkg/transaction/service"
	transactionstore "github.com/commonpool/backend/pkg/transaction/store"
	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	_ "net/http/pprof"
	"os"
	"sync"
	"testing"
)

var Db *gorm.DB
var AmqpClient mq.Client
var ResourceStore resourcestore.ResourceStore
var AuthStore store.UserStore
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
var UserService *service2.UserService
var ResourceService *service.ResourceService
var AuthHandler *authhandler.AuthHandler
var SessionHandler *session.Handler
var ChatHandler *chathandler.Handler
var GroupHandler *grouphandler.Handler
var ResourceHandler *resourcehandler.ResourceHandler
var UserHandler *authhandler.UserHandler
var TradingHandler *tradinghandler.TradingHandler

func TestMain(m *testing.M) {

	println("running main")

	appConfig, err := config.GetAppConfig(os.LookupEnv, ioutil.ReadFile)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	AmqpClient, err = mq.NewRabbitMqClient(ctx, appConfig.AmqpUrl)
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
	AuthStore = *store.NewUserStore(Driver)
	ChatStore = *chatstore.NewChatStore(Db)
	TradingStore = *tradingstore.NewTradingStore(Driver)
	GroupStore = *groupstore.NewGroupStore(Driver)
	ChatService = *chatservice.NewChatService(&AuthStore, AmqpClient, &ChatStore)
	GroupService = *groupservice.NewGroupService(&GroupStore, AmqpClient, ChatService, &AuthStore)
	TradingService = *tradingservice.NewTradingService(TradingStore, &ResourceStore, &AuthStore, ChatService, GroupService, TransactionService)
	UserService = service2.NewUserService(&AuthStore)
	ResourceService = service.NewResourceService(&ResourceStore)

	AuthHandler = authhandler.NewAuthHandler(Authorizer)
	//authHandler.Register(v1)

	SessionHandler = session.NewHandler(Authorizer)
	//sessionHandler.Register(v1)

	ChatHandler = chathandler.NewHandler(ChatService, TradingService, appConfig, Authorizer)
	//chatHandler.Register(v1)

	GroupHandler = grouphandler.NewHandler(GroupService, UserService, Authorizer)
	//groupHandler.Register(v1)

	ResourceHandler = resourcehandler.NewHandler(ResourceService, GroupService, UserService, Authorizer)
	//resourceHandler.Register(v1)

	UserHandler = authhandler.NewAuthHandler(UserService, Authorizer)
	//userHandler.Register(v1)

	TradingHandler = tradinghandler.NewTradingHandler(TradingService, GroupService, UserService, Authorizer)
	//userHandler.Register(v1)

	db.AutoMigrate(Db)

	cleanDb()

	session := Driver.GetSession()
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

func NewUser() *models.UserSession {
	userIncrementerMu.Lock()
	defer func() {
		userIncrementerMu.Unlock()
	}()
	userIncrementer++
	var userId = uuid.NewV4().String()
	userEmail := fmt.Sprintf("user%d@email.com", userIncrementer)
	userName := fmt.Sprintf("user%d", userIncrementer)
	return &models.UserSession{
		Username:        userName,
		Subject:         userId,
		Email:           userEmail,
		IsAuthenticated: true,
	}
}

var createUserLock sync.Mutex

func cleanDb() {

	session := Driver.GetSession()

	_, err := session.Run(`MATCH (n) DETACH DELETE n`, map[string]interface{}{})
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
