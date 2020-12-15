package resource

import (
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
)

type Service interface {
	GetResourcesByKeys(resourceKeys *resourcemodel.ResourceKeys) (*resourcemodel.Resources, error)
}
