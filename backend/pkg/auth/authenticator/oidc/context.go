package oidc

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/authenticator"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/exceptions"
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

// GetLoggedInUser gets the user session from the context
func WithContextUser(ctx context.Context, userKey, userName, email string) context.Context {
	ctx = context.WithValue(ctx, authenticator.SubjectKey, userKey)
	ctx = context.WithValue(ctx, authenticator.SubjectUsernameKey, userName)
	ctx = context.WithValue(ctx, authenticator.SubjectEmailKey, email)
	ctx = context.WithValue(ctx, authenticator.IsAuthenticatedKey, true)
	return ctx
}
