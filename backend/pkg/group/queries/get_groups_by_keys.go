package queries

import (
	"fmt"
	"github.com/commonpool/backend/pkg/group/readmodels"
	"github.com/commonpool/backend/pkg/keys"
	"gorm.io/gorm"
)

type GetGroupByKeys struct {
	db *gorm.DB
}

func NewGetGroupByKeys(db *gorm.DB) *GetGroupByKeys {
	return &GetGroupByKeys{
		db: db,
	}
}

func (q *GetGroupByKeys) Get(groupKeys *keys.GroupKeys) ([]*readmodels.GroupReadModel, error) {

	if len(groupKeys.Items) == 0 {
		return []*readmodels.GroupReadModel{}, nil
	}

	var rms []*readmodels.GroupReadModel

	sql := "group_key in ("
	var params []interface{}
	for i, groupKey := range groupKeys.Items {
		sql += "?"
		if i < len(groupKeys.Items)-1 {
			sql += ","
		}
		params = append(params, groupKey.String())
	}
	sql += ")"

	qry := q.db.Where(sql, params...).Find(&rms)
	if qry.Error != nil {
		return nil, qry.Error
	}

	if int(qry.RowsAffected) != len(groupKeys.Items) {
		return nil, fmt.Errorf("some groups were not found")
	}

	return rms, nil

}
