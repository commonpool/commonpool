package auth

import (
	"fmt"
	usermodel "github.com/commonpool/backend/pkg/user/model"
)

type UserNames map[usermodel.UserKey]string

func (u *UserNames) GetName(userKey usermodel.UserKey) (string, error) {
	userName, ok := (*u)[userKey]
	if !ok {
		return "", fmt.Errorf("name for '" + userKey.String() + "' not found")
	}
	return userName, nil
}
