package integration

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/config"
	"github.com/commonpool/backend/graph"
	"github.com/commonpool/backend/handler"
	"github.com/commonpool/backend/mock"
	"github.com/commonpool/backend/service"
	"github.com/commonpool/backend/store"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"sync"
	"testing"
)

var authenticatedUser = &auth.UserSession{}
var a *handler.Handler

var Db *gorm.DB
var AmqpClient amqp.Client
var ResourceStore store.ResourceStore
var AuthStore store.AuthStore
var ChatStore store.ChatStore
var TradingStore store.TradingStore
var GroupStore store.GroupStore
var ChatService service.ChatService
var TradingService service.TradingService
var GroupService service.GroupService
var Authorizer *mock.Authorizer
var Driver *graph.Neo4jGraphDriver

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

	Db = store.NewTestDb()
	ResourceStore = *store.NewResourceStore(Driver)
	AuthStore = *store.NewAuthStore(Db, Driver)
	ChatStore = *store.NewChatStore(Db, &AuthStore, AmqpClient)
	TradingStore = *store.NewTradingStore(Driver)
	GroupStore = *store.NewGroupStore(Driver)
	ChatService = *service.NewChatService(&AuthStore, &GroupStore, &ResourceStore, AmqpClient, &ChatStore)
	GroupService = *service.NewGroupService(&GroupStore, AmqpClient, ChatService, &AuthStore)
	TradingService = *service.NewTradingService(TradingStore, &ResourceStore, &AuthStore, ChatService, GroupService)

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
