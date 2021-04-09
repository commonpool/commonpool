package utils

import (
	"cp/pkg/api"
	"gorm.io/gorm"
	"strings"
)

func CountMessages(db *gorm.DB, posts []*api.Post) error {

	if len(posts) == 0 {
		return nil
	}

	var postMap = map[string]*api.Post{}
	var sb strings.Builder
	sb.WriteString("select thread_id, count(*) as count from messages where thread_id in (")
	var params []interface{}
	for i, post := range posts {
		postMap[post.ID] = post
		params = append(params, post.ID)
		sb.WriteString("?")
		if i < len(posts)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString(") group by thread_id")

	type resultStruct struct {
		ThreadID string
		Count    int
	}
	var result []*resultStruct
	if err := db.Raw(sb.String(), params...).Find(&result).Error; err != nil {
		return err
	}

	for _, result := range result {
		if post, ok := postMap[result.ThreadID]; ok {
			post.MessageCount = result.Count
		}
	}

	return nil

}

