package model

import "fmt"

type UserNames map[UserKey]string

func (u *UserNames) GetName(userKey UserKey) (string, error) {
	userName, ok := (*u)[userKey]
	if !ok {
		return "", fmt.Errorf("name for '" + userKey.String() + "' not found")
	}
	return userName, nil
}
