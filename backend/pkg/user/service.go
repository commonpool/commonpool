package user

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
)

type Service interface {
	GetUser(userKey keys.UserKey) (*User, error)
	GetUsername(userKey keys.UserKey) (string, error)
	Find(query Query) (*Users, error)
	GetByKeys(ctx context.Context, userKeys *keys.UserKeys) (*Users, error)
}
