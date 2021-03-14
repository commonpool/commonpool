package oidc

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/auth/domain"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/labstack/echo/v4"
)

// saveUserInfo saves the logged in user info th database
func SaveUserInfo(ctx context.Context, userRepo domain.UserRepository, sub string, email string, username string) error {

	user, err := userRepo.Load(ctx, keys.NewUserKey(sub))
	if err != nil {
		return err
	}
	userInfo := domain.UserInfo{
		Email:    email,
		Username: username,
	}
	if user.GetVersion() == 0 {
		if err := user.DiscoverUser(userInfo); err != nil {
			return err
		}
	} else {
		if err := user.ChangeUserInfo(userInfo); err != nil {
			return err
		}
	}
	if err := userRepo.Save(ctx, user); err != nil {
		return err
	}

	return nil
}

// SaveAuthenticatedUser when user logs in, update the context with the user info,
// and also saves the newly gotten user info in the db
func SaveAuthenticatedUser(c echo.Context, ctx context.Context, userRepo domain.UserRepository, sub string, username string, email string) error {
	SetAuthenticatedUser(c, sub, username, email)
	return SaveUserInfo(ctx, userRepo, sub, email, username)
}

// SaveAuthenticatedUser when user logs in, update the context with the user info,
// and also saves the newly gotten user info in the db
func SetAuthenticatedUser(c echo.Context, sub string, username string, email string) {
	SetIsAuthenticated(c, true)
	c.Set(authenticator.SubjectUsernameKey, username)
	c.Set(authenticator.SubjectEmailKey, email)
	c.Set(authenticator.SubjectKey, sub)
}

// SetIsAuthenticated marks the current user as authenticated
func SetIsAuthenticated(e echo.Context, isAuthenticated bool) {
	e.Set(authenticator.IsAuthenticatedKey, isAuthenticated)
}
