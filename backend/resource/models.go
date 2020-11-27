package resource

import (
	"fmt"
	"github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/model"
	"github.com/satori/go.uuid"
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
	key model.ResourceKey,
	resourceType ResourceType,
	createdBy string,
	summary string,
	description string,
	valueInHoursFrom int,
	valueInHoursTo int,
) Resource {
	return Resource{
		ID:               key.ID,
		Summary:          summary,
		Description:      description,
		CreatedBy:        createdBy,
		Type:             resourceType,
		ValueInHoursFrom: valueInHoursFrom,
		ValueInHoursTo:   valueInHoursTo,
	}
}

func (r *Resource) GetKey() model.ResourceKey {
	return model.NewResourceKey(r.ID)
}

func (r *Resource) GetOwnerKey() model.UserKey {
	return model.NewUserKey(r.CreatedBy)
}

type Resources struct {
	ItemMap map[model.ResourceKey]Resource
	Items   []Resource
}

func NewResources(items []Resource) *Resources {
	rsMap := map[model.ResourceKey]Resource{}
	for _, item := range items {
		rsMap[item.GetKey()] = item
	}
	return &Resources{
		Items:   items,
		ItemMap: rsMap,
	}
}

func (r *Resources) GetResource(key model.ResourceKey) (Resource, error) {
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

func (r *Resources) ContainsKey(key model.ResourceKey) bool {
	_, ok := r.ItemMap[key]
	return ok
}

type ResourceSharing struct {
	ResourceID uuid.UUID `gorm:"type:uuid;primary_key"`
	GroupID    uuid.UUID `gorm:"type:uuid;primary_key"`
}

func NewResourceSharing(resourceKey model.ResourceKey, groupKey model.GroupKey) ResourceSharing {
	return ResourceSharing{
		ResourceID: resourceKey.ID,
		GroupID:    groupKey.ID,
	}
}

type ResourceSharings struct {
	sharings map[model.ResourceKey][]ResourceSharing
}

func NewResourceSharings(sharings []ResourceSharing) (*ResourceSharings, error) {
	var result = map[model.ResourceKey][]ResourceSharing{}
	for _, sharing := range sharings {
		resourceKey, err := model.ParseResourceKey(sharing.ResourceID.String())
		if err != nil {
			return &ResourceSharings{}, err
		}
		_, ok := result[*resourceKey]
		if !ok {
			result[*resourceKey] = []ResourceSharing{}
		}
		result[*resourceKey] = append(result[*resourceKey], sharing)
	}
	return &ResourceSharings{sharings: result}, nil
}

func (s *ResourceSharings) GetAllGroupKeys() []model.GroupKey {
	groupMap := map[model.GroupKey]bool{}
	var groupKeys []model.GroupKey
	for _, sharing := range s.Items() {
		groupKey := model.NewGroupKey(sharing.GroupID)
		if !groupMap[groupKey] {
			groupMap[groupKey] = true
			groupKeys = append(groupKeys, groupKey)
		}
	}
	return groupKeys
}

func (s *ResourceSharings) GetSharingsForResource(key model.ResourceKey) *ResourceSharings {
	list, ok := s.sharings[key]
	if !ok {
		response, _ := NewResourceSharings([]ResourceSharing{})
		return response
	}
	response, _ := NewResourceSharings(list)
	return response
}

func (s *ResourceSharings) Items() []ResourceSharing {
	var result []ResourceSharing
	for _, sharings := range s.sharings {
		for _, sharing := range sharings {
			result = append(result, sharing)
		}
	}
	return result
}

type ResourceType int

const (
	ResourceOffer ResourceType = iota
	ResourceRequest
)

func ParseResourceType(s string) (*ResourceType, error) {
	var res ResourceType
	if s == "" {
		return nil, nil
	}
	if s == "0" {
		res = ResourceOffer
		return &res, nil
	}
	if s == "1" {
		res = ResourceRequest
		return &res, nil
	}

	err := errors.ErrParseResourceType(s)
	return nil, &err
}
