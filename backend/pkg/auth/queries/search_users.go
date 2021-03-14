package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/readmodel"
	"github.com/commonpool/backend/pkg/keys"
	"gorm.io/gorm"
)

type SearchUsers struct {
	db *gorm.DB
}

func NewSearchUsers(db *gorm.DB) *SearchUsers {
	return &SearchUsers{db: db}
}

type Query struct {
	Query      string
	Skip       int
	Take       int
	NotInGroup *keys.GroupKey
}

func (q *SearchUsers) Get(ctx context.Context, query Query) ([]*readmodel.UserReadModel, error) {

	qry := q.db.Model(&readmodel.UserReadModel{}).
		Offset(query.Skip).
		Limit(query.Take).
		Where("username like", "%"+query.Query+"%")

	if query.NotInGroup != nil {
		qry = qry.Where("not exists (select null from membership_read_models m where m.user_key = user_read_models.user_key and m.group_key = ?", query.NotInGroup.String())
	}

	var result []*readmodel.UserReadModel
	if err := qry.Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil

}
