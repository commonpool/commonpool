package user

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
)

type Store interface {
	GetByKey(userKey keys.UserKey) (*User, error)
	GetByKeys(ctx context.Context, userKeys *keys.UserKeys) (*Users, error)
	Upsert(userKey keys.UserKey, email string, username string) error
	GetUsername(userKey keys.UserKey) (string, error)
	Find(query Query) (*Users, error)
}

type Query struct {
	Query      string
	Skip       int
	Take       int
	NotInGroup *keys.GroupKey
}
