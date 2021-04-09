package utils

import (
	"cp/pkg/api"
	"gorm.io/gorm"
	"strings"
)

func PopulateTargets(db *gorm.DB, targets []*api.Target) error {

	var groupIds []string
	var userIds []string

	var groupMap = map[string]*api.Group{}
	var userMap = map[string]*api.User{}

	for _, target := range targets {
		if target.Type == api.GroupTarget {
			groupIds = append(groupIds, *target.GroupID)
		} else if target.Type == api.UserTarget {
			userIds = append(userIds, *target.UserID)
		}
	}

	if len(userIds) > 0 {
		var sb strings.Builder
		sb.WriteString("id in (")
		var params []interface{}
		for i, id := range userIds {
			params = append(params, id)
			sb.WriteString("?")
			if i < len(userIds)-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteString(")")
		var users []*api.User
		if err := db.Model(&api.User{}).Where(sb.String(), params...).Find(&users).Error; err != nil {
			return err
		}
		for _, user := range users {
			userMap[user.ID] = user
		}
	}

	if len(groupIds) > 0 {
		var sb strings.Builder
		sb.WriteString("id in (")
		var params []interface{}
		for i, id := range groupIds {
			params = append(params, id)
			sb.WriteString("?")
			if i < len(groupIds)-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteString(")")
		var groups []*api.Group
		if err := db.Model(&api.Group{}).Where(sb.String(), params...).Find(&groups).Error; err != nil {
			return err
		}
		for _, group := range groups {
			groupMap[group.ID] = group
		}
	}

	for _, target := range targets {
		if target.Type == api.GroupTarget {
			target.Group = groupMap[*target.GroupID]
		} else if target.Type == api.UserTarget {
			target.User = userMap[*target.UserID]
		}
	}

	return nil
}
