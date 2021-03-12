package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"gorm.io/gorm"
)

type GetUserAdministeredGroupKeys struct {
	db *gorm.DB
}

func NewGetUserAdministeredGroupKeys(db *gorm.DB) *GetUserAdministeredGroupKeys {
	return &GetUserAdministeredGroupKeys{db: db}
}

type administeredGroups struct {
	GroupKey keys.GroupKey
}

func (q *GetUserAdministeredGroupKeys) Get(ctx context.Context, userKey keys.UserKey) (*keys.GroupKeys, error) {
	var adminKeys []*administeredGroups
	err := q.db.Raw("select group_key from offer_user_membership_read_models where user_key = ? and (is_admin = true or is_owner = true)", userKey).
		Find(&adminKeys).
		Error
	if err != nil {
		return nil, err
	}
	var result []keys.GroupKey
	for _, key := range adminKeys {
		result = append(result, key.GroupKey)
	}
	return keys.NewGroupKeys(result), nil
}
