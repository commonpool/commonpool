package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/group/readmodels"
	"github.com/commonpool/backend/pkg/keys"
	"gorm.io/gorm"
)

type GetGroupMemberships struct {
	db *gorm.DB
}

func NewGetGroupMemberships(db *gorm.DB) *GetGroupMemberships {
	return &GetGroupMemberships{
		db: db,
	}
}

func (q *GetGroupMemberships) Get(ctx context.Context, groupKey keys.GroupKey, status *domain.MembershipStatus) ([]*readmodels.MembershipReadModel, error) {

	var rms []*readmodels.MembershipReadModel

	var sql = "group_key = ?"
	var params = []interface{}{groupKey.String()}

	if status != nil {
		sql = sql + " AND status = ?"
		params = append(params, status)
	}

	if err := q.db.Model(&readmodels.MembershipReadModel{}).Where(sql, params...).Find(&rms).Error; err != nil {
		return nil, err
	}

	return rms, nil

}
