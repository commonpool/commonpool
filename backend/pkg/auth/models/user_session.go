package models

import (
	"github.com/commonpool/backend/pkg/keys"
)

// UserSession Holds data for the currently authenticated user
type UserSession struct {
	Username        string
	Subject         string
	Email           string
	IsAuthenticated bool
}

var _ UserReference = &UserSession{}

// GetUserKey Gets the userKey from the UserSession
func (s *UserSession) GetUserKey() keys.UserKey {
	return keys.NewUserKey(s.Subject)
}

func (s *UserSession) Target() *keys.Target {
	return keys.NewUserKey(s.Subject).Target()
}

// GetUsername Gets the userName from the UserSession
func (s *UserSession) GetUsername() string {
	return s.Username
}
