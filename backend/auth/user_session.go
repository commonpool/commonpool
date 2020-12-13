package auth

import (
	"context"
	"github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/model"
	"github.com/labstack/echo/v4"
)

const (
	SubjectKey         = "auth-subject"
	SubjectUsernameKey = "auth-preferred-username"
	SubjectEmailKey    = "auth-email"
	IsAuthenticatedKey = "is-authenticated"
)

// UserSession Holds data for the currently authenticated user
type UserSession struct {
	Username        string
	Subject         string
	Email           string
	IsAuthenticated bool
}

var _ model.UserReference = &UserSession{}

// GetUserKey Gets the userKey from the UserSession
func (s *UserSession) GetUserKey() model.UserKey {
	return model.NewUserKey(s.Subject)
}

// GetUsername Gets the userName from the UserSession
func (s *UserSession) GetUsername() string {
	return s.Username
}

// saveAuthenticatedUser when user logs in, update the context with the user info,
// and also saves the newly gotten user info in the db
func saveAuthenticatedUser(c echo.Context, store Store, sub string, username string, email string) error {
	SetIsAuthenticated(c, true)
	c.Set(SubjectUsernameKey, username)
	c.Set(SubjectEmailKey, email)
	c.Set(SubjectKey, sub)
	return saveUserInfo(store, sub, email, username)
}

// saveUserInfo saves the logged in user info th database
func saveUserInfo(as Store, sub string, email string, username string) error {
	return as.Upsert(model.NewUserKey(sub), email, username)
}

// SetIsAuthenticated marks the current user as authenticated
func SetIsAuthenticated(e echo.Context, isAuthenticated bool) {
	e.Set(IsAuthenticatedKey, isAuthenticated)
}

// GetLoggedInUser gets the user session from the context
func GetLoggedInUser(ctx context.Context) (*UserSession, error) {

	valIntf := ctx.Value(IsAuthenticatedKey)

	if valIntf == nil {
		return nil, errors.ErrUnauthorized
	}

	if !valIntf.(bool) {
		return nil, errors.ErrUnauthorized
	}
	return &UserSession{
		Username:        ctx.Value(SubjectUsernameKey).(string),
		Subject:         ctx.Value(SubjectKey).(string),
		Email:           ctx.Value(SubjectEmailKey).(string),
		IsAuthenticated: true,
	}, nil

}
