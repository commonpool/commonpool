package models

import (
	"github.com/commonpool/backend/pkg/keys"
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

var _ UserReference = &User{}

func (u *User) GetUserKey() keys.UserKey {
	return keys.NewUserKey(u.ID)
}

func (u *User) GetUsername() string {
	return u.Username
}
