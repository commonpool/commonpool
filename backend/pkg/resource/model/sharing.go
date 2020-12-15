package model

import (
	groupmodel "github.com/commonpool/backend/pkg/group/model"
)

type Sharing struct {
	ResourceKey ResourceKey
	GroupKey    groupmodel.GroupKey
}

func NewResourceSharing(resourceKey ResourceKey, groupKey groupmodel.GroupKey) Sharing {
	return Sharing{
		ResourceKey: resourceKey,
		GroupKey:    groupKey,
	}
}
