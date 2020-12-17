package user

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
)

type Store interface {
	GetByKey(key keys.UserKey) (*User, error)
	GetByKeys(ctx context.Context, keys *keys.UserKeys) (*Users, error)
	Upsert(key keys.UserKey, email string, username string) error
	GetUsername(key keys.UserKey) (string, error)
	Find(query Query) (*Users, error)
}

type Query struct {
	Query      string
	Skip       int
	Take       int
	NotInGroup *keys.GroupKey
}
