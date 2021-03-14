package queries

import (
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/auth/readmodel"
	"github.com/commonpool/backend/pkg/keys"
	"gorm.io/gorm"
	"strings"
)

type GetUsersByKeys struct {
	db *gorm.DB
}

func NewGetUserByKeys(db *gorm.DB) *GetUsersByKeys {
	return &GetUsersByKeys{
		db: db,
	}
}

func (q *GetUsersByKeys) Get(userKeys *keys.UserKeys) (map[keys.UserKey]*models.User, error) {
	var users []*readmodel.UserReadModel
	var result = map[keys.UserKey]*models.User{}
	if userKeys.IsEmpty() {
		return result, nil
	}

	var sb strings.Builder
	var params []interface{}
	sb.WriteString("user_key in (")
	for i, userKey := range userKeys.Items {
		params = append(params, userKey.String())
		sb.WriteString("?")
		if i < len(userKeys.Items)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString(")")

	qry := q.db.Model(&readmodel.UserReadModel{}).Where(sb.String(), params...).Find(&users)
	if qry.Error != nil {
		return nil, qry.Error
	}

	for _, user := range users {
		result[user.UserKey] = &models.User{
			ID:       user.UserKey.String(),
			Username: user.Username,
			Email:    user.Email,
		}
	}

	return result, nil
}
