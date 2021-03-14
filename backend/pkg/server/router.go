package server

import (
	"github.com/commonpool/backend/pkg/validation"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"sync"
)

func NewRouter() *echo.Echo {
	e := echo.New()

	e.GET("/debug/pprof/", fromHandlerFunc(pprof.Index).Handle)
	e.GET("/debug/pprof/heap", fromHTTPHandler(pprof.Handler("heap")).Handle)
	e.GET("/debug/pprof/goroutine", fromHTTPHandler(pprof.Handler("goroutine")).Handle)
	e.GET("/debug/pprof/block", fromHTTPHandler(pprof.Handler("block")).Handle)
	e.GET("/debug/pprof/threadcreate", fromHTTPHandler(pprof.Handler("threadcreate")).Handle)
	e.GET("/debug/pprof/cmdline", fromHandlerFunc(pprof.Cmdline).Handle)
	e.GET("/debug/pprof/profile", fromHandlerFunc(pprof.Profile).Handle)
	e.GET("/debug/pprof/symbol", fromHandlerFunc(pprof.Symbol).Handle)

	e.Logger.SetLevel(log.DEBUG)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	e.Validator = validation.DefaultValidator
	return e
}

func fromHandlerFunc(serveHTTP func(w http.ResponseWriter, r *http.Request)) *customEchoHandler {
	return &customEchoHandler{httpHandler: &customHTTPHandler{serveHTTP: serveHTTP}}
}

func fromHTTPHandler(httpHandler http.Handler) *customEchoHandler {
	return &customEchoHandler{httpHandler: httpHandler}
}

type customEchoHandler struct {
	httpHandler       http.Handler
	wrappedHandleFunc echo.HandlerFunc
	once              sync.Once
}

func (ceh *customEchoHandler) Handle(c echo.Context) error {
	ceh.once.Do(func() {
		ceh.wrappedHandleFunc = ceh.mustWrapHandleFunc(c)
	})
	return ceh.wrappedHandleFunc(c)
}

func (ceh *customEchoHandler) mustWrapHandleFunc(c echo.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		ceh.httpHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

type customHTTPHandler struct {
	serveHTTP func(w http.ResponseWriter, r *http.Request)
}

func (c *customHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.serveHTTP(w, r)
}
