package queries

import (
	"github.com/commonpool/backend/pkg/auth/readmodel"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/keys"
	"gorm.io/gorm"
)

type GetUser struct {
	db *gorm.DB
}

func NewGetUser(db *gorm.DB) *GetUser {
	return &GetUser{
		db: db,
	}
}

func (q *GetUser) Get(userKey keys.UserKey) (*readmodel.UserReadModel, error) {
	var user readmodel.UserReadModel
	qry := q.db.Model(&readmodel.UserReadModel{}).Where("user_key = ?", userKey.String()).Find(&user)

	if qry.Error != nil {
		return nil, qry.Error
	}
	if qry.RowsAffected == 0 {
		return nil, exceptions.ErrUserNotFound
	}
	return &user, nil
}
