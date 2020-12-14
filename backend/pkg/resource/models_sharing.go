package resource

import "github.com/commonpool/backend/model"

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
