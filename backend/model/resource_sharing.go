package model

import (
	uuid "github.com/satori/go.uuid"
)

type ResourceSharing struct {
	ResourceID uuid.UUID `gorm:"type:uuid;primary_key"`
	GroupID    uuid.UUID `gorm:"type:uuid;primary_key"`
}

func NewResourceSharing(resourceKey ResourceKey, groupKey GroupKey) ResourceSharing {
	return ResourceSharing{
		ResourceID: resourceKey.uuid,
		GroupID:    groupKey.ID,
	}
}

type ResourceSharingKey struct {
	GroupKey    GroupKey
	ResourceKey ResourceKey
}

func NewResourceSharingKey(resourceKey ResourceKey, groupKey GroupKey) ResourceSharingKey {
	return ResourceSharingKey{
		GroupKey:    groupKey,
		ResourceKey: resourceKey,
	}
}

type ResourceSharings struct {
	sharings map[ResourceKey][]ResourceSharing
}

func NewResourceSharings(sharings []ResourceSharing) (*ResourceSharings, error) {
	var result = map[ResourceKey][]ResourceSharing{}
	for _, sharing := range sharings {
		resourceKey, err := ParseResourceKey(sharing.ResourceID.String())
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

func (s *ResourceSharings) GetAllGroupKeys() []GroupKey {
	groupMap := map[GroupKey]bool{}
	var groupKeys []GroupKey
	for _, sharing := range s.Items() {
		groupKey := NewGroupKey(sharing.GroupID)
		if !groupMap[groupKey] {
			groupMap[groupKey] = true
			groupKeys = append(groupKeys, groupKey)
		}
	}
	return groupKeys
}

func (s *ResourceSharings) GetSharingsForResource(key ResourceKey) *ResourceSharings {
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
