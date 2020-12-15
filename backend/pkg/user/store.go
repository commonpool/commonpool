package user

import (
	"context"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	model2 "github.com/commonpool/backend/pkg/user/model"
	usermodel "github.com/commonpool/backend/pkg/user/model"
)

type Store interface {
	GetByKey(key usermodel.UserKey) (*model2.User, error)
	GetByKeys(ctx context.Context, keys []usermodel.UserKey) (*Users, error)
	Upsert(key usermodel.UserKey, email string, username string) error
	GetUsername(key usermodel.UserKey) (string, error)
	Find(query Query) ([]*model2.User, error)
}

type Query struct {
	Query      string
	Skip       int
	Take       int
	NotInGroup *groupmodel.GroupKey
}
