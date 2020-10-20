package store

import (
	"errors"
	"github.com/commonpool/backend/auth"
	errs "github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/model"
	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func (rs *UserStore) GetByKeys(keys []model.UserKey, r []*model.User) error {
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

func (rs *UserStore) Upsert(key model.UserKey, email string, username string) error {
	usr := &model.User{}
	err := rs.GetByKey(key, usr)
	if err != nil && errs.IsNotFoundError(err) {
		usr.Username = username
		usr.ID = key.String()
		usr.Email = email
		return rs.db.Create(usr).Error
	} else if err != nil {
		return err
	} else {
		usr.Username = username
		usr.Email = email
		return rs.db.Save(usr).Error
	}
}

var _ auth.Store = &UserStore{}

// GetByKey Gets a resource by keys
func (rs *UserStore) GetByKey(key model.UserKey, r *model.User) error {
	if err := rs.db.First(r, "id = ?", key.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewUserNotFoundError(key.String())
		}
		return err
	}
	return nil
}
