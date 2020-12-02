package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/commonpool/backend/auth"
	errs "github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/model"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"strings"
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

func (as *AuthStore) GetByKeys(ctx context.Context, keys []model.UserKey) (auth.Users, error) {

	var qryStrs []string
	var qryParams []interface{}

	if keys == nil || len(keys) == 0 {
		return auth.NewUsers([]auth.User{}), nil
	}

	for _, userKey := range keys {
		qryStrs = append(qryStrs, "?")
		qryParams = append(qryParams, userKey.String())
	}

	sqlWhere := "id IN (" + strings.Join(qryStrs, ",") + ")"

	var users []auth.User
	err := as.db.Model(auth.User{}).Where(sqlWhere, qryParams...).Find(&users).Error

	if err != nil {
		log.Error(err, "GetByKeys: could not get users by keys")
		return auth.Users{}, err
	}

	return auth.NewUsers(users), nil

}

func (as *AuthStore) Upsert(key model.UserKey, email string, username string) error {
	usr := &auth.User{}
	err := as.GetByKey(key, usr)

	if err != nil && errs.IsNotFoundError(err) {
		log.Info(err, "Upsert: user not found. Creating...")
		usr.Username = username
		usr.ID = key.String()
		usr.Email = email
		return as.db.Create(usr).Error
	} else if err != nil {
		log.Error(err, "Upsert: error while upserting user")
		return err
	} else {
		if usr.Username == username && usr.Email == email {
			return nil
		}

		updates := map[string]interface{}{
			"username": username,
			"email":    email,
		}
		err := as.db.Model(auth.User{}).Where("id = ?", key.String()).Updates(updates).Error
		if err != nil {
			return fmt.Errorf("could not upsert user: %s", err.Error())
		}
		return nil
	}

}

func (as *AuthStore) GetByKey(key model.UserKey, r *auth.User) error {
	if err := as.db.Model(&auth.User{}).First(r, "id = ?", key.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn(err, "GetByKey: user with id "+key.String()+" not found")
			response := errs.ErrUserNotFound(key.String())
			return &response
		}
		log.Error(err, "GetByKey: could not get user by key")
		return err
	}
	return nil
}

func (as *AuthStore) GetUsername(key model.UserKey) (string, error) {
	var user auth.User
	err := as.GetByKey(key, &user)
	if err != nil {
		log.Error(err, "GetUsername: could not get username")
		return "", err
	}
	return user.Username, err
}

func (as *AuthStore) Find(query auth.UserQuery) ([]auth.User, error) {
	var users []auth.User
	chain := as.db.Order("username asc")
	if query.Query != "" {
		chain = chain.Where("username like ?", query.Query+"%")
	}

	if query.NotInGroup != nil {
		chain = chain.
			Joins("LEFT OUTER JOIN memberships ON (memberships.user_id = users.id and memberships.group_id = ?)", query.NotInGroup.ID.String()).
			Where("group_id IS NULL")
	}

	err := chain.Offset(query.Skip).Limit(query.Take).Find(&users).Error
	if err != nil {
		log.Error(err, "Find: could not find users")
		return nil, err
	}

	return users, err
}
