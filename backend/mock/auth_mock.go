package mock

import (
	"context"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Authorizer struct {
	IsAuthorized       bool
	MockCurrentSession func() auth.UserSession
}

func (a *Authorizer) GetRedirectResponse(request *http.Request) (*auth.RedirectResponse, error) {
	panic("implement me")
}

var _ auth.IAuth = &Authorizer{}

func (a *Authorizer) GetAuthUserSession(c echo.Context) auth.UserSession {
	var isAuthenticated = c.Get(auth.IsAuthenticatedKey).(bool)
	if !isAuthenticated {
		return auth.UserSession{
			Username:        "",
			Subject:         "",
			Email:           "",
			IsAuthenticated: false,
		}
	}
	return auth.UserSession{
		Username:        c.Get(auth.SubjectUsernameKey).(string),
		Subject:         c.Get(auth.SubjectKey).(string),
		Email:           c.Get(auth.SubjectEmailKey).(string),
		IsAuthenticated: isAuthenticated,
	}
}

func (a *Authorizer) GetAuthUserKey(c echo.Context) model.UserKey {
	session := a.GetAuthUserSession(c)
	return session.GetUserKey()
}

func (a *Authorizer) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	}
}

func (a *Authorizer) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	}
}

func (a *Authorizer) Authenticate(redirectOnError bool) echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return handlerFunc(c)
		}
	}
}

func (a *Authorizer) GetAuthUserSession2(ctx context.Context) auth.UserSession {
	return a.MockCurrentSession()
}

func NewTestAuthorizer() *Authorizer {
	return &Authorizer{}
}
