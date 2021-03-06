package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/group/readmodels"
	"github.com/commonpool/backend/pkg/keys"
	"gorm.io/gorm"
)

type GetUserMemberships struct {
	db *gorm.DB
}

func NewGetUserMemberships(db *gorm.DB) *GetUserMemberships {
	return &GetUserMemberships{
		db: db,
	}
}

func (q *GetUserMemberships) Get(ctx context.Context, userKey keys.UserKey, status *domain.MembershipStatus) ([]*readmodels.MembershipReadModel, error) {

	var rms []*readmodels.MembershipReadModel

	var sql = "user_key = ?"
	var params = []interface{}{userKey.String()}

	if status != nil {
		sql = sql + " AND status = ?"
		params = append(params, status)
	}

	if err := q.db.Model(&readmodels.MembershipReadModel{}).Where(sql, params...).Find(&rms).Error; err != nil {
		return nil, err
	}

	return rms, nil

}
