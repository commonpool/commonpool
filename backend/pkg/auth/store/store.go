package store

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/keys"
)

type Store interface {
	GetByKey(userKey keys.UserKey) (*models.User, error)
	GetByKeys(ctx context.Context, userKeys *keys.UserKeys) (*models.Users, error)
	Upsert(userKey keys.UserKey, email string, username string) error
	GetUsername(userKey keys.UserKey) (string, error)
	Find(query Query) (*models.Users, error)
}

type Query struct {
	Query      string
	Skip       int
	Take       int
	NotInGroup *keys.GroupKey
}
