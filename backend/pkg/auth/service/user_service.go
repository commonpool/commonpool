package service

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/auth/store"
	"github.com/commonpool/backend/pkg/keys"
)

type UserService struct {
	userStore store.Store
}

func (u UserService) GetByKeys(ctx context.Context, keys *keys.UserKeys) (*models.Users, error) {
	return u.userStore.GetByKeys(ctx, keys)
}

func NewUserService(userStore store.Store) *UserService {
	return &UserService{
		userStore: userStore,
	}
}

func (u UserService) GetUser(key keys.UserKey) (*models.User, error) {
	return u.userStore.GetByKey(key)
}

func (u UserService) GetUsername(key keys.UserKey) (string, error) {
	return u.userStore.GetUsername(key)
}

func (u UserService) Find(query store.Query) (*models.Users, error) {
	return u.userStore.Find(query)
}

var _ Service = &UserService{}
