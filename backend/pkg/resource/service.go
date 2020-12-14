package resource

import "github.com/commonpool/backend/model"

type Service interface {
	GetResourcesByKeys(resourceKeys *model.ResourceKeys) (*Resources, error)
}
