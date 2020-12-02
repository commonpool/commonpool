package integration

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/config"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/handler"
	"github.com/commonpool/backend/mock"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/service"
	"github.com/commonpool/backend/store"
	"github.com/commonpool/backend/trading"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
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

func TestMain(m *testing.M) {

	println("running main")

	ctx := context.Background()

	AmqpClient, _ = amqp.NewRabbitMqClient(ctx, "amqp://guest:guest@192.168.39.47:31991/")
	Authorizer = mock.NewTestAuthorizer()
	Authorizer.MockCurrentSession = func() auth.UserSession {
		if authenticatedUser == nil {
			return auth.UserSession{
				IsAuthenticated: false,
			}
		}
		return *authenticatedUser
	}

	appConfig := &config.AppConfig{}

	Db = store.NewTestDb()
	ResourceStore = *store.NewResourceStore(Db)
	AuthStore = *store.NewAuthStore(Db)
	ChatStore = *store.NewChatStore(Db, &AuthStore, AmqpClient)
	TradingStore = *store.NewTradingStore(Db)
	GroupStore = *store.NewGroupStore(Db, AmqpClient)
	ChatService = *service.NewChatService(&AuthStore, &GroupStore, &ResourceStore, AmqpClient, &ChatStore)
	TradingService = *service.NewTradingService(TradingStore, &ResourceStore, &AuthStore, ChatService)
	GroupService = *service.NewGroupService(&GroupStore, AmqpClient, ChatService, &AuthStore)

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

	Db.Delete(resource.Resource{}, "1 = 1")
	Db.Delete(resource.Sharing{}, "1 = 1")
	Db.Delete(trading.Offer{}, "1 = 1")
	Db.Delete(trading.OfferItem{}, "1 = 1")
	Db.Delete(trading.OfferDecision{}, "1 = 1")
	Db.Delete(chat.Channel{}, "1 = 1")
	Db.Delete(chat.ChannelSubscription{}, "1 = 1")
	Db.Delete(store.Message{}, "1 = 1")
	Db.Delete(group.Group{}, "1 = 1")
	Db.Delete(group.Membership{}, "1 = 1")
}
