package handler

import (
	"github.com/commonpool/backend/pkg/group/readmodels"
	"time"
)

type Group struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func NewGroup(group *readmodels.GroupReadModel) *Group {
	return &Group{
		ID:          group.GroupKey,
		CreatedAt:   group.CreatedAt,
		Name:        group.Name,
		Description: group.Description,
	}
}
