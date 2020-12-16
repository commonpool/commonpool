package model

import (
	"github.com/commonpool/backend/pkg/group"
)

type ResourceSharingKey struct {
	GroupKey    group.GroupKey
	ResourceKey ResourceKey
}

//goland:noinspection GoUnusedExportedFunction
func NewResourceSharingKey(resourceKey ResourceKey, groupKey group.GroupKey) ResourceSharingKey {
	return ResourceSharingKey{
		GroupKey:    groupKey,
		ResourceKey: resourceKey,
	}
}
