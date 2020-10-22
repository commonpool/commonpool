package store

import (
	"github.com/commonpool/backend/model"
	"gorm.io/gorm"
)

type MessageStore struct {
	db *gorm.DB
}

func (rs *UserStore) GetLatestThreads(keys []model.UserKey, r []*model.User) error {
	for i, key := range keys {
		usr := &model.User{}
		err := rs.GetByKey(key, usr)
		if err != nil {
			return err
		}
		r[i] = usr
	}
	return nil
}

func (rs *UserStore) GetThreadMessages(keys []model.UserKey, r []*model.User) error {
	for i, key := range keys {
		usr := &model.User{}
		err := rs.GetByKey(key, usr)
		if err != nil {
			return err
		}
		r[i] = usr
	}
	return nil
}

