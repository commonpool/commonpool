package auth

import (
	"context"
	"github.com/commonpool/backend/pkg/exceptions"
	exceptions2 "github.com/commonpool/backend/pkg/user"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
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

var _ usermodel.UserReference = &UserSession{}

// GetUserKey Gets the userKey from the UserSession
func (s *UserSession) GetUserKey() usermodel.UserKey {
	return usermodel.NewUserKey(s.Subject)
}

// GetUsername Gets the userName from the UserSession
func (s *UserSession) GetUsername() string {
	return s.Username
}

func SetAuthenticatedUser(c echo.Context, username, subject, email string) {
	c.Set(IsAuthenticatedKey, true)
	c.Set(SubjectUsernameKey, username)
	c.Set(SubjectEmailKey, email)
	c.Set(SubjectKey, subject)
}

func SetContextAuthenticatedUser(c context.Context, username, subject, email string) context.Context {
	c = context.WithValue(c, IsAuthenticatedKey, true)
	c = context.WithValue(c, SubjectUsernameKey, username)
	c = context.WithValue(c, SubjectEmailKey, email)
	c = context.WithValue(c, SubjectKey, subject)
	return c
}

// saveAuthenticatedUser when user logs in, update the context with the user info,
// and also saves the newly gotten user info in the db
func saveAuthenticatedUser(c echo.Context, store exceptions2.Store, sub string, username string, email string) error {
	SetIsAuthenticated(c, true)
	c.Set(SubjectUsernameKey, username)
	c.Set(SubjectEmailKey, email)
	c.Set(SubjectKey, sub)
	return saveUserInfo(store, sub, email, username)
}

// saveUserInfo saves the logged in user info th database
func saveUserInfo(as exceptions2.Store, sub string, email string, username string) error {
	return as.Upsert(usermodel.NewUserKey(sub), email, username)
}

// SetIsAuthenticated marks the current user as authenticated
func SetIsAuthenticated(e echo.Context, isAuthenticated bool) {
	e.Set(IsAuthenticatedKey, isAuthenticated)
}

// GetLoggedInUser gets the user session from the context
func GetLoggedInUser(ctx context.Context) (*UserSession, error) {

	valIntf := ctx.Value(IsAuthenticatedKey)

	if valIntf == nil {
		return nil, exceptions.ErrUnauthorized
	}

	if !valIntf.(bool) {
		return nil, exceptions.ErrUnauthorized
	}
	return &UserSession{
		Username:        ctx.Value(SubjectUsernameKey).(string),
		Subject:         ctx.Value(SubjectKey).(string),
		Email:           ctx.Value(SubjectEmailKey).(string),
		IsAuthenticated: true,
	}, nil

}
