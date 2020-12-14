package handler

import (
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/db"
	group2 "github.com/commonpool/backend/pkg/group"
	resource2 "github.com/commonpool/backend/pkg/resource"
	trading2 "github.com/commonpool/backend/pkg/trading"
	"github.com/commonpool/backend/pkg/user"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

var (
	d          *gorm.DB
	rs         resource2.Store
	as         user.Store
	cs         chat.Store
	ts         trading2.Store
	gs         group2.Store
	h          *Handler
	e          *echo.Echo
	userSub1   = "user-1-sub"
	userSub2   = "user-2-sub"
	username1  = "user1"
	username2  = "user2"
	user1Email = "user1@example.com"
	user2Email = "user2@example.com"
	user1      = &auth.UserSession{
		Username:        username1,
		Subject:         userSub1,
		Email:           user1Email,
		IsAuthenticated: true,
	}
	user2 = &auth.UserSession{
		Username:        username2,
		Subject:         userSub2,
		Email:           user2Email,
		IsAuthenticated: true,
	}
	authenticatedUser = user1
	users             = []*auth.UserSession{user1, user2}
)

func mockLoggedInAs(user *auth.UserSession) {
	authenticatedUser = user
}

/**type Handler struct {
	amqp           amqp.AmqpClient
	resourceStore  resource.Store
	authStore      auth.Store
	authorization  auth.Authenticator
	chatStore      chat.Store
	tradingStore   trading.Store
	groupService   group.Service
	config         config.AppConfig
	chatService    chat.Service
	tradingService trading.Service
}

/**
resourceStore = store.NewResourceStore(db)
	authStore = store.NewAuthStore(db)
	chatStore := store.NewChatStore(db, authStore, amqpCli)
	tradingStore := store.NewTradingStore(db)
	groupStore := store.NewGroupStore(db, amqpCli)

	chatService := service.NewChatService(authStore, groupStore, resourceStore, amqpCli, chatStore)
	tradingService := service.NewTradingService(tradingStore, resourceStore, authStore, chatService)
	groupService := service.NewGroupService(groupStore, amqpCli, chatService, authStore)
*/

func setup() {

	//
	// fakeServer := amqpServer.NewServer("amqp://localhost:5672/%2f")
	// err := fakeServer.Start()
	// if err != nil {
	// 	panic(err)
	// }
	//
	// amqpClient, err := NewWabbitMqClient(fakeServer)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// // Setup and migrate database
	// d = store.NewTestDb()
	// store.AutoMigrate(d)
	//
	// // Create the different stores
	// rs = store.NewResourceStore(d)
	// as = store.NewAuthStore(d)
	// cs = store.NewChatStore(d)
	// ts = store.NewTradingStore(d)
	// gs := store.NewGroupStore(d)
	//
	// // Mock authorization
	// authorizer := auth.NewTestAuthorizer()
	// authorizer.MockCurrentSession = func() auth.UserSession {
	// 	if authenticatedUser == nil {
	// 		return auth.UserSession{
	// 			IsAuthenticated: false,
	// 		}
	// 	}
	// 	return *authenticatedUser
	// }
	//
	// // Create handler
	// h = NewHandler(rs, as, cs, ts, gs, authorizer, amqpClient, config.AppConfig{})
	//
	// // Create users
	// for _, user := range users {
	// 	userKey := model.NewUserKey(user.Subject)
	// 	err := as.Upsert(userKey, user.Email, user.Username)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	//
	// e = router.NewRouter()
}

func tearDown() {
	if err := db.DropTestDB(); err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}
