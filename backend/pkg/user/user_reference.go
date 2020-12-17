package user

import "github.com/commonpool/backend/pkg/keys"

type UserReference interface {
	GetUserKey() keys.UserKey
	GetUsername() string
}
