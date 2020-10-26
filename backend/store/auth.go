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

func (as *AuthStore) GetByKeys(keys []model.UserKey, r []*model.User) error {
	for i, key := range keys {
		usr := &model.User{}
		err := as.GetByKey(key, usr)
		if err != nil {
			return err
		}
		r[i] = usr
	}
	return nil
}

func (as *AuthStore) Upsert(key model.UserKey, email string, username string) error {
	usr := &model.User{}
	err := as.GetByKey(key, usr)
	if err != nil && errs.IsNotFoundError(err) {
		usr.Username = username
		usr.ID = key.String()
		usr.Email = email
		return as.db.Create(usr).Error
	} else if err != nil {
		return err
	} else {
		usr.Username = username
		usr.Email = email
		return as.db.Save(usr).Error
	}
}

func (as *AuthStore) GetByKey(key model.UserKey, r *model.User) error {
	if err := as.db.First(r, "id = ?", key.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := errs.ErrUserNotFound(key.String())
			return &response
		}
		return err
	}
	return nil
}

func (as *AuthStore) GetUsername(key model.UserKey) (string, error) {
	var user model.User
	err := as.GetByKey(key, &user)
	if err != nil {
		return "", err
	}
	return user.Username, err
}

func (as *AuthStore) Find(query auth.UserQuery) ([]model.User, error) {
	var users []model.User

	chain := as.db

	if query.Query != "" {
		chain = chain.Where("username like ?", query.Query+"%")
	}

	err := chain.Offset(query.Skip).Limit(query.Take).Find(&users).Error
	return users, err
}
