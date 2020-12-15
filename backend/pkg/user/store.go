package user

import (
	"context"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	usermodel "github.com/commonpool/backend/pkg/user/model"
)

type Store interface {
	GetByKey(key usermodel.UserKey) (*usermodel.User, error)
	GetByKeys(ctx context.Context, keys *usermodel.UserKeys) (*Users, error)
	Upsert(key usermodel.UserKey, email string, username string) error
	GetUsername(key usermodel.UserKey) (string, error)
	Find(query Query) (*Users, error)
}

type Query struct {
	Query      string
	Skip       int
	Take       int
	NotInGroup *groupmodel.GroupKey
}
