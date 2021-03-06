package domain

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
)

type UserRepository interface {
	Load(ctx context.Context, userKey keys.UserKey) (*User, error)
	Save(ctx context.Context, user *User) error
}
