package oidc

import (
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/auth/store"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/labstack/echo/v4"
)

// saveUserInfo saves the logged in user info th database
func SaveUserInfo(as store.Store, sub string, email string, username string) error {
	return as.Upsert(keys.NewUserKey(sub), email, username)
}

// SaveAuthenticatedUser when user logs in, update the context with the user info,
// and also saves the newly gotten user info in the db
func SaveAuthenticatedUser(c echo.Context, store store.Store, sub string, username string, email string) error {
	SetIsAuthenticated(c, true)
	c.Set(authenticator.SubjectUsernameKey, username)
	c.Set(authenticator.SubjectEmailKey, email)
	c.Set(authenticator.SubjectKey, sub)
	return SaveUserInfo(store, sub, email, username)
}

// SetIsAuthenticated marks the current user as authenticated
func SetIsAuthenticated(e echo.Context, isAuthenticated bool) {
	e.Set(authenticator.IsAuthenticatedKey, isAuthenticated)
}
