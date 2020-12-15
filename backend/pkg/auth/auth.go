package auth

import (
	"context"
	usermodel "github.com/commonpool/backend/pkg/user/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Interface for authorization module
type Authenticator interface {
	Login() echo.HandlerFunc
	Logout() echo.HandlerFunc
	Authenticate(redirectOnError bool) echo.MiddlewareFunc
	GetRedirectResponse(request *http.Request) (*RedirectResponse, error)
	GetLoggedInUser(ctx context.Context) (usermodel.UserReference, error)
}

// Ascertain that OidcAuthenticator implements Authenticator
var _ Authenticator = &OidcAuthenticator{}
