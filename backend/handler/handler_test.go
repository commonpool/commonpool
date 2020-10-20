package handler

import (
	"github.com/commonpool/backend/auth"
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
	d  *gorm.DB
	rs resource.Store
	as auth.Store
	h  *Handler
	e  *echo.Echo
)

var authenticatedUser auth.UserSession = struct {
	Username        string
	Subject         string
	IsAuthenticated bool
}{Username: "user", Subject: "subject", IsAuthenticated: true}

func setup() {
	d = store.NewTestDb()
	store.AutoMigrate(d)
	rs = store.NewResourceStore(d)
	as = store.NewAuthStore(d)
	authorizer := auth.NewTestAuthorizer()
	authorizer.MockAuthenticatedUser = func() auth.UserSession {
		return authenticatedUser
	}

	h = NewHandler(rs, as, authorizer)
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
