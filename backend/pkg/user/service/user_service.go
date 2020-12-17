package service

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/user"
)

type UserService struct {
	userStore user.Store
}

func (u UserService) GetByKeys(ctx context.Context, keys *keys.UserKeys) (*user.Users, error) {
	return u.userStore.GetByKeys(ctx, keys)
}

func NewUserService(userStore user.Store) *UserService {
	return &UserService{
		userStore: userStore,
	}
}

func (u UserService) GetUser(key keys.UserKey) (*user.User, error) {
	return u.userStore.GetByKey(key)
}

func (u UserService) GetUsername(key keys.UserKey) (string, error) {
	return u.userStore.GetUsername(key)
}

func (u UserService) Find(query user.Query) (*user.Users, error) {
	return u.userStore.Find(query)
}

var _ user.Service = &UserService{}
