package integration

import (
	"context"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/config"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/handler"
	"github.com/commonpool/backend/mock"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/service"
	"github.com/commonpool/backend/store"
	"github.com/commonpool/backend/trading"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"os"
	"testing"
)

var authenticatedUser = &auth.UserSession{}
var a *handler.Handler

var Db *gorm.DB
var AmqpClient amqp.AmqpClient
var ResourceStore store.ResourceStore
var AuthStore store.AuthStore
var ChatStore store.ChatStore
var TradingStore store.TradingStore
var GroupStore store.GroupStore
var ChatService service.ChatService
var TradingService service.TradingService
var GroupService service.GroupService
var User1KeyStr = uuid.NewV4().String()
var User1Key = model.NewUserKey(User1KeyStr)
var User1 *auth.UserSession
var User2KeyStr = uuid.NewV4().String()
var User2Key = model.NewUserKey(User1KeyStr)
var User2 *auth.UserSession
var User3KeyStr = uuid.NewV4().String()
var User3Key = model.NewUserKey(User1KeyStr)
var User3 *auth.UserSession
var Authorizer *mock.MockAuthorizer

func TestMain(m *testing.M) {

	println("running main")

	User1 = &auth.UserSession{
		Username:        "user1",
		Subject:         User1KeyStr,
		Email:           "user1@email.com",
		IsAuthenticated: true,
	}

	User2 = &auth.UserSession{
		Username:        "user2",
		Subject:         User2KeyStr,
		Email:           "user2@email.com",
		IsAuthenticated: true,
	}

	User3 = &auth.UserSession{
		Username:        "user3",
		Subject:         User3KeyStr,
		Email:           "user3@email.com",
		IsAuthenticated: true,
	}

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

	os.Exit(m.Run())

}

func teardown() {

	ctx := context.Background()

	Db.Delete(auth.User{}, "1 = 1")
	Db.Delete(resource.Resource{}, "1 = 1")
	Db.Delete(resource.ResourceSharing{}, "1 = 1")
	Db.Delete(trading.Offer{}, "1 = 1")
	Db.Delete(trading.OfferItem{}, "1 = 1")
	Db.Delete(trading.OfferDecision{}, "1 = 1")
	Db.Delete(chat.Channel{}, "1 = 1")
	Db.Delete(chat.ChannelSubscription{}, "1 = 1")
	Db.Delete(store.Message{}, "1 = 1")
	Db.Delete(group.Group{}, "1 = 1")
	Db.Delete(group.Membership{}, "1 = 1")

	ch, _ := AmqpClient.GetChannel()
	_ = ch.ExchangeDelete(ctx, User1Key.GetExchangeName(), false, false)
	_ = ch.ExchangeDelete(ctx, User2Key.GetExchangeName(), false, false)
	_ = ch.ExchangeDelete(ctx, User3Key.GetExchangeName(), false, false)
}

func setup() {

	PanicIfError(AuthStore.Upsert(User1.GetUserKey(), User1.Email, User1.Username))
	PanicIfError(AuthStore.Upsert(User2.GetUserKey(), User2.Email, User2.Username))
	PanicIfError(AuthStore.Upsert(User3.GetUserKey(), User3.Email, User3.Username))

}
