package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/group/readmodels"
	"github.com/commonpool/backend/pkg/keys"
	"gorm.io/gorm"
)

type GetMembershipReadModel struct {
	db *gorm.DB
}

func NewGetMembership(db *gorm.DB) *GetMembershipReadModel {
	return &GetMembershipReadModel{
		db: db,
	}
}

func (q *GetMembershipReadModel) Get(ctx context.Context, membershipKey keys.MembershipKey) (*readmodels.MembershipReadModel, error) {
	var rm readmodels.MembershipReadModel

	query := q.db.Model(&readmodels.MembershipReadModel{}).
		Where("group_key = ? and user_key = ?", membershipKey.GroupKey.String(), membershipKey.UserKey.String()).
		Find(&rm)

	if err := query.Error; err != nil {
		return nil, err
	}

	if query.RowsAffected == 0 {
		return nil, exceptions.ErrMembershipNotFound
	}

	return &rm, nil
}
