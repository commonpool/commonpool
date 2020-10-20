package auth

import "github.com/commonpool/backend/model"

type Store interface {
	GetByKey(key model.UserKey, r *model.User) error
	GetByKeys(keys []model.UserKey, r []*model.User) error
	Upsert(key model.UserKey, email string, username string) error
}
