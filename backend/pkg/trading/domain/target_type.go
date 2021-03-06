package domain

import "github.com/commonpool/backend/pkg/exceptions"

type TargetType string

const (
	UserTarget  TargetType = "user"
	GroupTarget TargetType = "group"
)

func (a TargetType) IsGroup() bool {
	return a == GroupTarget
}

func (a TargetType) IsUser() bool {
	return a == UserTarget
}

func ParseOfferItemTargetType(str string) (TargetType, error) {
	if str == "user" {
		return UserTarget, nil
	} else if str == "group" {
		return GroupTarget, nil
	} else {
		return "", exceptions.ErrInvalidTargetType
	}
}
