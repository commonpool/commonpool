package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/readmodel"
	"github.com/commonpool/backend/pkg/keys"
	"gorm.io/gorm"
)

type GetUsersForGroupInvite struct {
	db *gorm.DB
}

func NewGetUsersForGroupInvite(db *gorm.DB) *GetUsersForGroupInvite {
	return &GetUsersForGroupInvite{db: db}
}

func (q *GetUsersForGroupInvite) Get(ctx context.Context, groupKey keys.GroupKey, query string, skip, take int) ([]*readmodel.UserReadModel, error) {
	var result []*readmodel.UserReadModel
	query += "%"

	err := q.db.Raw(`
		select *
		from user_read_models 
		where
			username like ?
			and not exists(
				select null 
				from membership_read_models 
				where group_key = ? 
				and user_read_models.user_key = membership_read_models.user_key)`,
		query,
		groupKey.String()).Find(&result).Error

	return result, err
}
