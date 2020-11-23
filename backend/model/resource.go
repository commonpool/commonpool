package model

import (
	"fmt"
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
	return NewResourceKey(r.ID)
}

func (r *Resource) GetUserKey() UserKey {
	return NewUserKey(r.CreatedBy)
}

type Resources struct {
	ItemMap map[ResourceKey]Resource
	Items   []Resource
}

func NewResources(items []Resource) *Resources {
	rsMap := map[ResourceKey]Resource{}
	for _, item := range items {
		rsMap[item.GetKey()] = item
	}
	return &Resources{
		Items:   items,
		ItemMap: rsMap,
	}
}

func (r *Resources) GetResource(key ResourceKey) (Resource, error) {
	rs, ok := r.ItemMap[key]
	if !ok {
		return Resource{}, fmt.Errorf("resource not found")
	}
	return rs, nil
}

func (r *Resources) Append(resource Resource) *Resources {
	items := append(r.Items, resource)
	return NewResources(items)
}

func (r *Resources) Contains(resource Resource) bool {
	return r.ContainsKey(resource.GetKey())
}

func (r *Resources) ContainsKey(key ResourceKey) bool {
	_, ok := r.ItemMap[key]
	return ok
}
