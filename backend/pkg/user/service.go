package user

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
)

type Service interface {
	GetUser(key keys.UserKey) (*User, error)
	GetUsername(key keys.UserKey) (string, error)
	Find(query Query) (*Users, error)
	GetByKeys(ctx context.Context, keys *keys.UserKeys) (*Users, error)
}
