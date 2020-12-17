package auth

import (
	"fmt"
	"github.com/commonpool/backend/pkg/keys"
)

type UserNames map[keys.UserKey]string

func (u *UserNames) GetName(userKey keys.UserKey) (string, error) {
	userName, ok := (*u)[userKey]
	if !ok {
		return "", fmt.Errorf("name for '" + userKey.String() + "' not found")
	}
	return userName, nil
}
