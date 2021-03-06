package authenticator

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Interface for authorization module
type Authenticator interface {
	Login() echo.HandlerFunc
	Logout() echo.HandlerFunc
	Authenticate(redirectOnError bool) echo.MiddlewareFunc
	GetRedirectResponse(request *http.Request) (*RedirectResponse, error)
	GetLoggedInUser(ctx context.Context) (models.UserReference, error)
	Register(c *echo.Group)
}
