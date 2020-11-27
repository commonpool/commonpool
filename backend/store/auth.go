package store

import (
	"context"
	"errors"
	"github.com/commonpool/backend/auth"
	errs "github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/utils"
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

	var result []auth.User
	usedKeys := map[model.UserKey]bool{}

	err := utils.Partition(len(keys), 999, func(i1 int, i2 int) error {

		var qryStrs []string
		var qryParams []interface{}

		for _, item := range keys[i1:i2] {
			if usedKeys[item] {
				continue
			}
			usedKeys[item] = true
			qryStrs = append(qryStrs, "?")
			qryParams = append(qryParams, item.String())
		}
		sqlWhere := "id IN (" + strings.Join(qryStrs, ",") + ")"

		var partitioned []auth.User
		err := as.db.Model(auth.User{}).Where(sqlWhere, qryParams...).Find(&partitioned).Error
		if err != nil {
			log.Error(err, "GetByKeys: could not get users by keys")
			return err
		}
		for _, user := range partitioned {
			result = append(result, user)
		}
		return nil
	})

	if err != nil {
		return auth.Users{}, err
	}

	return auth.NewUsers(result), nil
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
			log.Error(err, "could not upsert user")
			return err
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
