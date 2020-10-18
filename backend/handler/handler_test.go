package handler

import (
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
	h  *Handler
	e  *echo.Echo
)

func setup() {
	d = store.NewTestDb()
	store.AutoMigrate(d)
	rs = store.NewResourceStore(d)
	h = NewHandler(rs)
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
