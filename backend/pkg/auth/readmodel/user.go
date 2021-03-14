package readmodel

import "github.com/commonpool/backend/pkg/keys"

type UserReadModel struct {
	UserKey  keys.UserKey `gorm:"index:idx_user_read_model_uq,unique;primaryKey"`
	Email    string
	Username string
	Version  int `gorm:"index:idx_user_read_model_uq,unique"`
}
