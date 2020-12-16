package model

import (
	"github.com/commonpool/backend/pkg/group"
)

type Sharing struct {
	ResourceKey ResourceKey
	GroupKey    group.GroupKey
}

func NewResourceSharing(resourceKey ResourceKey, groupKey group.GroupKey) Sharing {
	return Sharing{
		ResourceKey: resourceKey,
		GroupKey:    groupKey,
	}
}
