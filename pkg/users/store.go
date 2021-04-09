package users

import (
	"cp/pkg/api"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type Store interface {
	Get(userID string) (*api.User, error)
	GetByKeys(userIDs []string) ([]*api.User, error)
	Upsert(user *api.User) error
	Save(user *api.User) error
}

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{db: db}
}

func (u UserStore) Get(userID string) (*api.User, error) {
	var result api.User
	if err := u.db.First(&result, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (u UserStore) GetByKeys(userIDs []string) ([]*api.User, error) {
	var result []*api.User

	if len(userIDs) == 0 {
		return []*api.User{}, nil
	}

	var sb strings.Builder
	var params []interface{}
	sb.WriteString("id in (")
	for i, userID := range userIDs {
		sb.WriteString("?")
		params = append(params, userID)
		if i < len(userIDs)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString(")")

	if err := u.db.Where(sb.String(), params...).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (u UserStore) Save(user *api.User) error {
	return u.db.Save(user).Error
}

func (u UserStore) Upsert(user *api.User) error {
	return u.db.Debug().Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"username":     user.Username,
			"email":        user.Email,
		}),
	}).Create(user).Error
}
