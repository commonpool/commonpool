package user

import (
	"github.com/commonpool/backend/model"
	"time"
)

type User struct {
	ID        string `gorm:"primary_key" mapstructure:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Username  string     `gorm:"not null" mapstructure:"username"`
	Email     string     `gorm:"not null" mapstructure:"email"`
}

var _ model.UserReference = &User{}

func (u *User) GetUserKey() model.UserKey {
	return model.NewUserKey(u.ID)
}

func (u *User) GetUsername() string {
	return u.Username
}
