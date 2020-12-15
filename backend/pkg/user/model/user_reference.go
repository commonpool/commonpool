package model

type UserReference interface {
	GetUserKey() UserKey
	GetUsername() string
}
