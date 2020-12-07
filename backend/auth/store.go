package auth

import (
	"context"
	"github.com/commonpool/backend/model"
)

type Store interface {
	GetByKey(key model.UserKey) (*User, error)
	GetByKeys(ctx context.Context, keys []model.UserKey) (*Users, error)
	Upsert(key model.UserKey, email string, username string) error
	GetUsername(key model.UserKey) (string, error)
	Find(query UserQuery) ([]*User, error)
}

type UserQuery struct {
	Query      string
	Skip       int
	Take       int
	NotInGroup *model.GroupKey
}
