package oidc

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/labstack/echo/v4"
)

// GetLoggedInUser gets the user session from the context
func GetLoggedInUser(ctx context.Context) (*models.UserSession, error) {

	valIntf := ctx.Value(authenticator.IsAuthenticatedKey)

	if valIntf == nil {
		return nil, exceptions.ErrUnauthorized
	}

	if !valIntf.(bool) {
		return nil, exceptions.ErrUnauthorized
	}
	return &models.UserSession{
		Username:        ctx.Value(authenticator.SubjectUsernameKey).(string),
		Subject:         ctx.Value(authenticator.SubjectKey).(string),
		Email:           ctx.Value(authenticator.SubjectEmailKey).(string),
		IsAuthenticated: true,
	}, nil

}

func SetAuthenticatedUser(c echo.Context, username, subject, email string) {
	c.Set(authenticator.IsAuthenticatedKey, true)
	c.Set(authenticator.SubjectUsernameKey, username)
	c.Set(authenticator.SubjectEmailKey, email)
	c.Set(authenticator.SubjectKey, subject)
}

func SetContextAuthenticatedUser(c context.Context, username, subject, email string) context.Context {
	c = context.WithValue(c, authenticator.IsAuthenticatedKey, true)
	c = context.WithValue(c, authenticator.SubjectUsernameKey, username)
	c = context.WithValue(c, authenticator.SubjectEmailKey, email)
	c = context.WithValue(c, authenticator.SubjectKey, subject)
	return c
}
