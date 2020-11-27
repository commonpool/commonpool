package auth

import (
	"context"
	"github.com/commonpool/backend/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

type MockAuthorizer struct {
	IsAuthorized       bool
	MockCurrentSession func() UserSession
}

var _ IAuth = &MockAuthorizer{}

func (a *MockAuthorizer) GetAuthUserSession(c echo.Context) UserSession {
	return a.MockCurrentSession()
}

func (a *MockAuthorizer) GetAuthUserKey(c echo.Context) model.UserKey {
	return model.NewUserKey(a.MockCurrentSession().Subject)
}

func (a *MockAuthorizer) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	}
}

func (a *MockAuthorizer) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	}
}

func (a *MockAuthorizer) Authenticate(redirectOnError bool) echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return handlerFunc(c)
		}
	}
}

func (a *MockAuthorizer) GetAuthUserSession2(ctx context.Context) UserSession {
	return a.MockCurrentSession()
}

func NewTestAuthorizer() *MockAuthorizer {
	return &MockAuthorizer{}
}
