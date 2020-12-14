package auth

import (
	"fmt"
	"github.com/commonpool/backend/model"
)

type UserNames map[model.UserKey]string

func (u *UserNames) GetName(userKey model.UserKey) (string, error) {
	userName, ok := (*u)[userKey]
	if !ok {
		return "", fmt.Errorf("name for '" + userKey.String() + "' not found")
	}
	return userName, nil
}
