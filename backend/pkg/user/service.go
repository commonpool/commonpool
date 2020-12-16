package user

import (
	"context"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
)

type Service interface {
	GetUser(key usermodel.UserKey) (*usermodel.User, error)
	GetUsername(key usermodel.UserKey) (string, error)
	Find(query Query) (*Users, error)
	GetByKeys(ctx context.Context, keys *usermodel.UserKeys) (*Users, error)
}
