package usermodel

type UserReference interface {
	GetUserKey() UserKey
	GetUsername() string
}
