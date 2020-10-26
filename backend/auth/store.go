package auth

import "github.com/commonpool/backend/model"

type Store interface {
	GetByKey(key model.UserKey, r *model.User) error
	GetByKeys(keys []model.UserKey, r []*model.User) error
	Upsert(key model.UserKey, email string, username string) error
	GetUsername(key model.UserKey) (string, error)
	Find(query UserQuery) ([]model.User, error)
}

type UserQuery struct {
	Query      string
	Skip       int
	Take       int
	NotInGroup *model.GroupKey
}
