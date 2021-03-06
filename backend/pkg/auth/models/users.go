package models

import (
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/keys"
)

type Users struct {
	Items   []*User
	userMap map[keys.UserKey]*User
}

func NewUsers(u []*User) *Users {
	var users []*User
	var userMap = map[keys.UserKey]*User{}
	for _, user := range u {
		userKey := user.GetUserKey()
		if _, ok := userMap[userKey]; ok {
			continue
		}
		users = append(users, user)
		userMap[userKey] = user
	}
	if users == nil {
		users = []*User{}
	}
	return &Users{
		Items:   users,
		userMap: userMap,
	}
}

func NewEmptyUsers() *Users {
	return &Users{
		Items:   []*User{},
		userMap: map[keys.UserKey]*User{},
	}
}

func (u *Users) GetUser(key keys.UserKey) (*User, error) {
	user, ok := u.userMap[key]
	if !ok {
		return nil, exceptions.ErrUserNotFound
	}
	return user, nil
}

func (u *Users) Contains(user keys.UserKey) bool {
	_, ok := u.userMap[user]
	return ok
}

func (u *Users) Append(user *User) *Users {
	if u.Contains(user.GetUserKey()) {
		return NewUsers(u.Items)
	}
	newItems := append(u.Items, user)
	return NewUsers(newItems)
}

func (u *Users) AppendAll(users *Users) *Users {
	var items = u.Items
	for _, user := range users.Items {
		items = append(items, user)
	}
	return NewUsers(items)
}

func (u *Users) GetUserKeys() *keys.UserKeys {
	var userKeys []keys.UserKey
	for _, item := range u.Items {
		userKeys = append(userKeys, item.GetUserKey())
	}
	return keys.NewUserKeys(userKeys)
}

func (u *Users) GetUserCount() int {
	return len(u.Items)
}
