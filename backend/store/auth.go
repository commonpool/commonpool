package store

import (
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

func (as *AuthStore) GetByKeys(keys []model.UserKey) ([]model.User, error) {

	var result []model.User
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

		var partitioned []model.User
		err := as.db.Model(model.User{}).Where(sqlWhere, qryParams...).Find(&partitioned).Error
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
		return nil, err
	}

	return result, nil
}

func (as *AuthStore) Upsert(key model.UserKey, email string, username string) error {
	usr := &model.User{}
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
		err := as.db.Model(model.User{}).Where("id = ?", key.String()).Updates(updates).Error
		if err != nil {
			log.Error(err, "could not upsert user")
			return err
		}
		return nil
	}

}

func (as *AuthStore) GetByKey(key model.UserKey, r *model.User) error {
	if err := as.db.Model(&model.User{}).First(r, "id = ?", key.String()).Error; err != nil {
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
	var user model.User
	err := as.GetByKey(key, &user)
	if err != nil {
		log.Error(err, "GetUsername: could not get username")
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
	if query.NotInGroup != nil {
		chain = chain.
			Joins("left join memberships on memberships.user_id = users.id ").
			Where("memberships.group_id != ? or memberships.group_id is null", query.NotInGroup.ID.String())
	}

	err := chain.Offset(query.Skip).Limit(query.Take).Find(&users).Error
	if err != nil {
		log.Error(err, "Find: could not find users")
	}

	return users, err
}
