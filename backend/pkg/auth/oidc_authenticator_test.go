package auth

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContextIsNotAuthenticatedAtAll(t *testing.T) {
	oidc := OidcAuthenticator{}
	ctx := context.Background()
	session := oidc.GetAuthUserSession(ctx)
	assert.False(t, session.IsAuthenticated)
}

func TestContextIsNotAuthenticated(t *testing.T) {
	oidc := OidcAuthenticator{}
	ctx := context.Background()
	ctx = context.WithValue(ctx, IsAuthenticatedKey, false)
	session := oidc.GetAuthUserSession(ctx)
	assert.False(t, session.IsAuthenticated)
}

func TestContextIsAuthenticated(t *testing.T) {
	oidc := OidcAuthenticator{}
	ctx := context.Background()
	ctx = context.WithValue(ctx, IsAuthenticatedKey, true)
	ctx = context.WithValue(ctx, SubjectKey, "user1")
	ctx = context.WithValue(ctx, SubjectEmailKey, "user1@email.com")
	ctx = context.WithValue(ctx, SubjectUsernameKey, "username")
	session := oidc.GetAuthUserSession(ctx)
	assert.True(t, session.IsAuthenticated)
	assert.Equal(t, "user1", session.Subject)
	assert.Equal(t, "user1@email.com", session.Email)
	assert.Equal(t, "username", session.Username)
}
