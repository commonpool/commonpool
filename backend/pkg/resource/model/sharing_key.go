package model

import groupmodel "github.com/commonpool/backend/pkg/group/model"

type ResourceSharingKey struct {
	GroupKey    groupmodel.GroupKey
	ResourceKey ResourceKey
}

//goland:noinspection GoUnusedExportedFunction
func NewResourceSharingKey(resourceKey ResourceKey, groupKey groupmodel.GroupKey) ResourceSharingKey {
	return ResourceSharingKey{
		GroupKey:    groupKey,
		ResourceKey: resourceKey,
	}
}
