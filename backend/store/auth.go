package store

import (
	"errors"
	"github.com/commonpool/backend/auth"
	errs "github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/model"
	"gorm.io/gorm"
)

type AuthStore struct {
	db *gorm.DB
}

var _ auth.Store = &AuthStore{}

func NewAuthStore(db *gorm.DB) *AuthStore {
	return &AuthStore{
		db: db,
	}
}

type UserStore struct {
	db *gorm.DB
}

func (rs *AuthStore) GetByKeys(keys []model.UserKey, r []*model.User) error {
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

func (rs *AuthStore) Upsert(key model.UserKey, email string, username string) error {
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

// GetByKey Gets a resource by keys
func (rs *AuthStore) GetByKey(key model.UserKey, r *model.User) error {
	if err := rs.db.First(r, "id = ?", key.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := errs.ErrUserNotFound(key.String())
			return &response
		}
		return err
	}
	return nil
}
