package main

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	_ "github.com/commonpool/backend/docs"
	"github.com/commonpool/backend/handler"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/router"
	"github.com/commonpool/backend/store"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

var (
	d  *gorm.DB
	rs resource.Store
	e  *echo.Echo
)


// @title resources api
// @version 1.0
// @description resources api
// @termsOfService http://swagger.io/terms
// @host 127.0.0.1:8585
// @basePath /api/v1
func main() {
	r := router.NewRouter()

	r.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := r.Group("/api/v1")

	d = store.NewTestDb()
	store.AutoMigrate(d)
	rs = store.NewResourceStore(d)
	h := handler.NewHandler(rs)
	h.Register(v1)
	r.Logger.Fatal(r.Start("127.0.0.1:8585"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
