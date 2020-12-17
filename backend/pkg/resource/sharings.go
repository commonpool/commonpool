package resource

import (
	"github.com/commonpool/backend/pkg/keys"
)

type Sharings struct {
	sharingMap map[keys.ResourceKey][]*Sharing
}

func (s *Sharings) AppendAll(appendSharings *Sharings) {
	for _, item := range appendSharings.Items() {
		if _, ok := s.sharingMap[item.ResourceKey]; !ok {
			s.sharingMap[item.ResourceKey] = []*Sharing{}
		}
		s.sharingMap[item.ResourceKey] = append(s.sharingMap[item.ResourceKey], &Sharing{
			ResourceKey: item.ResourceKey,
			GroupKey:    item.GroupKey,
		})
	}
}

func NewResourceSharings(sharings []*Sharing) *Sharings {
	var result = map[keys.ResourceKey][]*Sharing{}
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
		sharingMap: map[keys.ResourceKey][]*Sharing{},
	}
}

func (s *Sharings) GetAllGroupKeys() *keys.GroupKeys {
	groupMap := map[keys.GroupKey]bool{}
	var groupKeys []keys.GroupKey
	for _, sharing := range s.Items() {
		if !groupMap[sharing.GroupKey] {
			groupMap[sharing.GroupKey] = true
			groupKeys = append(groupKeys, sharing.GroupKey)
		}
	}
	return keys.NewGroupKeys(groupKeys)
}

func (s *Sharings) GetSharingsForResource(key keys.ResourceKey) *Sharings {
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
