package handler

import (
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/router"
	"github.com/commonpool/backend/store"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

var (
	d          *gorm.DB
	rs         resource.Store
	as         auth.Store
	cs         chat.Store
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

func setup() {
	d = store.NewTestDb()
	store.AutoMigrate(d)

	rs = store.NewResourceStore(d)
	as = store.NewAuthStore(d)
	cs = store.NewChatStore(d)

	authorizer := auth.NewTestAuthorizer()
	authorizer.MockCurrentSession = func() auth.UserSession {
		if authenticatedUser == nil {
			return auth.UserSession{
				IsAuthenticated: false,
			}
		}
		return *authenticatedUser
	}

	h = NewHandler(rs, as, cs, authorizer)

	for _, user := range users {
		userKey := model.NewUserKey(user.Subject)
		err := as.Upsert(userKey, user.Email, user.Username)
		if err != nil {
			panic(err)
		}
	}

	e = router.NewRouter()
}

func tearDown() {
	if err := store.DropTestDB(); err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}
