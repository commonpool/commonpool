package queries

import (
	"github.com/commonpool/backend/pkg/auth/readmodel"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/keys"
	"gorm.io/gorm"
)

type GetUsername struct {
	db *gorm.DB
}

func NewGetUsername(db *gorm.DB) *GetUsername {
	return &GetUsername{
		db: db,
	}
}

func (q *GetUsername) Get(userKey keys.UserKey) (string, error) {
	var rm readmodel.UserReadModel
	qry := q.db.Find(&rm, "user_key = ?", userKey.String())
	if qry.Error != nil {
		return "", qry.Error
	}
	if qry.RowsAffected == 0 {
		return "", exceptions.ErrUserNotFound
	}
	return rm.Username, nil
}
