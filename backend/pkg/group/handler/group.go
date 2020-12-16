package handler

import (
	"github.com/commonpool/backend/pkg/group"
	"time"
)

type Group struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func NewGroup(group *group.Group) *Group {
	return &Group{
		ID:          group.Key.String(),
		CreatedAt:   group.CreatedAt,
		Name:        group.Name,
		Description: group.Description,
	}
}
