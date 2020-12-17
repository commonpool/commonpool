package resource

import (
	"github.com/commonpool/backend/pkg/keys"
)

type Sharing struct {
	ResourceKey keys.ResourceKey
	GroupKey    keys.GroupKey
}

func NewResourceSharing(resourceKey keys.ResourceKey, groupKey keys.GroupKey) Sharing {
	return Sharing{
		ResourceKey: resourceKey,
		GroupKey:    groupKey,
	}
}
