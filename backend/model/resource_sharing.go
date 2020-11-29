package model

type ResourceSharingKey struct {
	GroupKey    GroupKey
	ResourceKey ResourceKey
}

//goland:noinspection GoUnusedExportedFunction
func NewResourceSharingKey(resourceKey ResourceKey, groupKey GroupKey) ResourceSharingKey {
	return ResourceSharingKey{
		GroupKey:    groupKey,
		ResourceKey: resourceKey,
	}
}
