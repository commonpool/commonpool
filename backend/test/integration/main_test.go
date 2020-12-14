package integration

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	chatservice "github.com/commonpool/backend/chat/service"
	chatstore "github.com/commonpool/backend/chat/store"
	"github.com/commonpool/backend/config"
	"github.com/commonpool/backend/graph"
	"github.com/commonpool/backend/handler"
	"github.com/commonpool/backend/mock"
	"github.com/commonpool/backend/pkg/chat"
	service2 "github.com/commonpool/backend/pkg/group/service"
	store4 "github.com/commonpool/backend/pkg/group/store"
	store5 "github.com/commonpool/backend/pkg/resource/store"
	service3 "github.com/commonpool/backend/pkg/trading/service"
	store2 "github.com/commonpool/backend/pkg/trading/store"
	service4 "github.com/commonpool/backend/pkg/transaction/service"
	store3 "github.com/commonpool/backend/pkg/transaction/store"
	store "github.com/commonpool/backend/store"
	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	_ "net/http/pprof"
	"os"
	"sync"
	"testing"
)

var authenticatedUser = &auth.UserSession{}
var a *handler.Handler

var Db *gorm.DB
var AmqpClient amqp.Client
var ResourceStore store5.ResourceStore
var AuthStore store.AuthStore
var ChatStore chatstore.ChatStore
var TradingStore store2.TradingStore
var GroupStore store4.GroupStore
var ChatService chatservice.ChatService
var TradingService service3.TradingService
var GroupService service2.GroupService
var Authorizer *mock.Authorizer
var Driver *graph.Neo4jGraphDriver
var TransactionStore *store3.TransactionStore
var TransactionService *service4.TransactionService

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
	Authorizer = mock.NewTestAuthorizer()
	Authorizer.MockCurrentSession = func() auth.UserSession {
		if authenticatedUser == nil {
			return auth.UserSession{
				IsAuthenticated: false,
			}
		}
		return *authenticatedUser
	}

	err = graph.InitGraphDatabase(nil, appConfig)
	if err != nil {
		panic(err)
	}

	Driver, err = graph.NewNeo4jDriver(appConfig, appConfig.Neo4jDatabase)
	if err != nil {
		panic(err)
	}

	Db = getDb(appConfig)

	TransactionStore = store3.NewTransactionStore(Db)
	TransactionService = service4.NewTransactionService(TransactionStore)
	ResourceStore = *store5.NewResourceStore(Driver, TransactionService)
	AuthStore = *store.NewAuthStore(Db, Driver)
	ChatStore = *chatstore.NewChatStore(Db, &AuthStore, AmqpClient)
	TradingStore = *store2.NewTradingStore(Driver)
	GroupStore = *store4.NewGroupStore(Driver)
	ChatService = *chatservice.NewChatService(&AuthStore, &GroupStore, &ResourceStore, AmqpClient, &ChatStore)
	GroupService = *service2.NewGroupService(&GroupStore, AmqpClient, ChatService, &AuthStore)
	TradingService = *service3.NewTradingService(TradingStore, &ResourceStore, &AuthStore, ChatService, GroupService, TransactionService)

	store.AutoMigrate(Db)

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
	Db.Delete(auth.User{}, "1 = 1")

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

	Db.Delete(chat.Channel{}, "1 = 1")
	Db.Delete(chat.ChannelSubscription{}, "1 = 1")
	Db.Delete(store.Message{}, "1 = 1")
}

func getDb(appConfig *config.AppConfig) *gorm.DB {
	cs := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", appConfig.DbHost, appConfig.DbUsername, appConfig.DbPassword, appConfig.DbName, appConfig.DbPort)
	db, err := gorm.Open(postgres.Open(cs), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
