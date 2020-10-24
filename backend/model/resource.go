package model

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Resource struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time   `sql:"index"`
	Summary          string       `gorm:"not null;"`
	Description      string       `gorm:"not null;"`
	CreatedBy        string       `gorm:"not null;"`
	Type             ResourceType `gorm:"not null;"`
	ValueInHoursFrom int          `gorm:"not null;'"`
	ValueInHoursTo   int          `gorm:"not null"`
}

func NewResource(
	key ResourceKey,
	resourceType ResourceType,
	createdBy string,
	summary string,
	description string,
	valueInHoursFrom int,
	valueInHoursTo int,
) Resource {
	return Resource{
		ID:               key.uuid,
		Summary:          summary,
		Description:      description,
		CreatedBy:        createdBy,
		Type:             resourceType,
		ValueInHoursFrom: valueInHoursFrom,
		ValueInHoursTo:   valueInHoursTo,
	}
}

func (r *Resource) GetKey() ResourceKey {
	return ResourceKey{
		uuid: r.ID,
	}
}

func (r *Resource) GetUserKey() UserKey {
	return NewUserKey(r.CreatedBy)
}
