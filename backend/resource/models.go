package resource

import (
	"fmt"
	"github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/model"
	"time"
)

type Resource struct {
	Key              model.ResourceKey
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
	Summary          string
	Description      string
	CreatedBy        string
	Type             Type
	ValueInHoursFrom int
	ValueInHoursTo   int
}

func NewResource(
	key model.ResourceKey,
	resourceType Type,
	createdBy string,
	summary string,
	description string,
	valueInHoursFrom int,
	valueInHoursTo int,
) Resource {
	return Resource{
		Key:              key,
		Summary:          summary,
		Description:      description,
		CreatedBy:        createdBy,
		Type:             resourceType,
		ValueInHoursFrom: valueInHoursFrom,
		ValueInHoursTo:   valueInHoursTo,
	}
}

func (r *Resource) GetKey() model.ResourceKey {
	return r.Key
}

func (r *Resource) GetOwnerKey() model.UserKey {
	return model.NewUserKey(r.CreatedBy)
}

type Resources struct {
	ItemMap map[model.ResourceKey]*Resource
	Items   []*Resource
}

func NewResources(items []*Resource) *Resources {
	rsMap := map[model.ResourceKey]*Resource{}
	for _, item := range items {
		rsMap[item.GetKey()] = item
	}
	return &Resources{
		Items:   items,
		ItemMap: rsMap,
	}
}

func NewEmptyResources() *Resources {
	return &Resources{
		ItemMap: map[model.ResourceKey]*Resource{},
		Items:   []*Resource{},
	}
}

func (r *Resources) GetResource(key model.ResourceKey) (*Resource, error) {
	rs, ok := r.ItemMap[key]
	if !ok {
		return nil, fmt.Errorf("resource not found")
	}
	return rs, nil
}

func (r *Resources) Append(resource *Resource) *Resources {
	items := append(r.Items, resource)
	return NewResources(items)
}

func (r *Resources) Contains(resource *Resource) bool {
	return r.ContainsKey(resource.GetKey())
}

func (r *Resources) ContainsKey(key model.ResourceKey) bool {
	_, ok := r.ItemMap[key]
	return ok
}

func (r *Resources) GetKeys() *model.ResourceKeys {
	var resourceKeys []model.ResourceKey
	for _, resource := range r.Items {
		resourceKeys = append(resourceKeys, resource.GetKey())
	}
	if resourceKeys == nil {
		resourceKeys = []model.ResourceKey{}
	}
	return model.NewResourceKeys(resourceKeys)
}

type Sharing struct {
	ResourceKey model.ResourceKey
	GroupKey    model.GroupKey
}

func NewResourceSharing(resourceKey model.ResourceKey, groupKey model.GroupKey) Sharing {
	return Sharing{
		ResourceKey: resourceKey,
		GroupKey:    groupKey,
	}
}

type Sharings struct {
	sharingMap map[model.ResourceKey][]*Sharing
}

func NewResourceSharings(sharings []*Sharing) *Sharings {
	var result = map[model.ResourceKey][]*Sharing{}
	for _, sharing := range sharings {
		_, ok := result[sharing.ResourceKey]
		if !ok {
			result[sharing.ResourceKey] = []*Sharing{}
		}
		result[sharing.ResourceKey] = append(result[sharing.ResourceKey], sharing)
	}
	return &Sharings{sharingMap: result}
}

func NewEmptyResourceSharings() *Sharings {
	return &Sharings{
		sharingMap: map[model.ResourceKey][]*Sharing{},
	}
}

func (s *Sharings) GetAllGroupKeys() *model.GroupKeys {
	groupMap := map[model.GroupKey]bool{}
	var groupKeys []model.GroupKey
	for _, sharing := range s.Items() {
		if !groupMap[sharing.GroupKey] {
			groupMap[sharing.GroupKey] = true
			groupKeys = append(groupKeys, sharing.GroupKey)
		}
	}
	return model.NewGroupKeys(groupKeys)
}

func (s *Sharings) GetSharingsForResource(key model.ResourceKey) *Sharings {
	list, ok := s.sharingMap[key]
	if !ok {
		return NewResourceSharings([]*Sharing{})
	}
	return NewResourceSharings(list)
}

func (s *Sharings) Items() []*Sharing {
	var result []*Sharing
	if s.sharingMap == nil {
		return []*Sharing{}
	}
	for _, sharingMapEntry := range s.sharingMap {
		for _, sharing := range sharingMapEntry {
			result = append(result, sharing)
		}
	}
	if result == nil {
		return []*Sharing{}
	}
	return result
}

type Type int

const (
	Offer Type = iota
	Request
)

func ParseResourceType(s string) (*Type, error) {
	var res Type
	if s == "" {
		return nil, nil
	}
	if s == "0" {
		res = Offer
		return &res, nil
	}
	if s == "1" {
		res = Request
		return &res, nil
	}

	err := errors.ErrParseResourceType(s)
	return nil, &err
}
