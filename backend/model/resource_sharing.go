package model

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
