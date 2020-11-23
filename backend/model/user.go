package model

import (
	"fmt"
	"time"
)

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

type Users struct {
	Items   []User
	userMap map[UserKey]User
}

func NewUsers(u []User) Users {
	var users []User
	var userMap = map[UserKey]User{}
	for _, user := range u {
		users = append(users, user)
		userMap[user.GetKey()] = user
	}
	return Users{
		Items:   users,
		userMap: userMap,
	}
}

func (u *Users) GetUser(key UserKey) (User,error) {
	user, ok := u.userMap[key]
	if !ok {
		return User{}, fmt.Errorf("user not found")
	}
	return user, nil
}

func (u *Users) Append(user User) Users {
	items := append(u.Items, user)
	return NewUsers(items)
}

func (u *Users) AppendAll(users Users) Users {
	items := u.Items
	for _, user := range users.Items {
		items = append(items, user)
	}
	return NewUsers(items)
}

func (u *Users) GetUserKeys() UserKeys {
	var userKeys []UserKey
	for _, item := range u.Items {
		userKeys = append(userKeys, item.GetKey())
	}
	return NewUserKeys(userKeys)
}


func (u *Users) GetUserCount() int {
	return len(u.Items)
}
