package model

import "time"

type User struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Username  string     `gorm:"not null"`
	Email     string     `gorm:"not null"`
}

func (u *User) GetKey() UserKey {
	return NewUserKey(u.ID)
}
