package user

import (
	"fmt"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
)

type Users struct {
	Items   []*usermodel.User
	userMap map[usermodel.UserKey]*usermodel.User
}

func NewUsers(u []*usermodel.User) *Users {
	var users []*usermodel.User
	var userMap = map[usermodel.UserKey]*usermodel.User{}
	for _, user := range u {
		userKey := user.GetUserKey()
		if _, ok := userMap[userKey]; ok {
			continue
		}
		users = append(users, user)
		userMap[userKey] = user
	}
	if users == nil {
		users = []*usermodel.User{}
	}
	return &Users{
		Items:   users,
		userMap: userMap,
	}
}

func NewEmptyUsers() *Users {
	return &Users{
		Items:   []*usermodel.User{},
		userMap: map[usermodel.UserKey]*usermodel.User{},
	}
}

func (u *Users) GetUser(key usermodel.UserKey) (*usermodel.User, error) {
	user, ok := u.userMap[key]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (u *Users) Contains(user usermodel.UserKey) bool {
	_, ok := u.userMap[user]
	return ok
}

func (u *Users) Append(user *usermodel.User) *Users {
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

func (u *Users) GetUserKeys() *usermodel.UserKeys {
	var userKeys []usermodel.UserKey
	for _, item := range u.Items {
		userKeys = append(userKeys, item.GetUserKey())
	}
	return usermodel.NewUserKeys(userKeys)
}

func (u *Users) GetUserCount() int {
	return len(u.Items)
}
