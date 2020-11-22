package model

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Group struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	CreatedBy   string
	CreatedAt   time.Time
	Name        string
	Description string
}

func (o *Group) GetKey() GroupKey {
	return NewGroupKey(o.ID)
}

type Groups struct {
	items map[GroupKey]Group
	order []int
}
