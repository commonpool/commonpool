package service

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/auth/store"
	"github.com/commonpool/backend/pkg/keys"
)

type Service interface {
	GetUser(userKey keys.UserKey) (*models.User, error)
	GetUsername(userKey keys.UserKey) (string, error)
	Find(query store.Query) (*models.Users, error)
	GetByKeys(ctx context.Context, userKeys *keys.UserKeys) (*models.Users, error)
}
