package auth

import (
	"context"
	"github.com/commonpool/backend/model"
	"github.com/labstack/echo/v4"
)

// Interface for authorization module
type IAuth interface {
	Login() echo.HandlerFunc
	Logout() echo.HandlerFunc
	Authenticate(redirectOnError bool) echo.MiddlewareFunc
	GetAuthUserSession(c echo.Context) UserSession
	GetAuthUserSession2(ctx context.Context) UserSession
	GetAuthUserKey(c echo.Context) model.UserKey
}

// Ascertain that OidcAuthenticator implements IAuth
var _ IAuth = &OidcAuthenticator{}
