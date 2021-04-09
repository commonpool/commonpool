package utils

import "cp/pkg/api"

func UserMap(users []*api.User) map[string]*api.User {
	var result = map[string]*api.User{}
	for _, user := range users {
		result[user.ID] = user
	}
	return result
}
