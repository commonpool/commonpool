package user

import (
	"context"
	"github.com/commonpool/backend/pkg/group"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
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
	NotInGroup *group.GroupKey
}
