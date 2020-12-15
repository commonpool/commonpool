package model

import (
	usermodel "github.com/commonpool/backend/pkg/user/model"
	"time"
)

type Resource struct {
	Key              ResourceKey
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
	Summary          string
	Description      string
	CreatedBy        string
	Type             Type
	ValueInHoursFrom int
	ValueInHoursTo   int
	SubType          SubType
}

func NewResource(
	key ResourceKey,
	resourceType Type,
	subType SubType,
	createdBy string,
	summary string,
	description string,
	valueInHoursFrom int,
	valueInHoursTo int) Resource {
	return Resource{
		Key:              key,
		Summary:          summary,
		Description:      description,
		CreatedBy:        createdBy,
		Type:             resourceType,
		SubType:          subType,
		ValueInHoursFrom: valueInHoursFrom,
		ValueInHoursTo:   valueInHoursTo,
	}
}

func (r *Resource) GetKey() ResourceKey {
	return r.Key
}

func (r *Resource) GetOwnerKey() usermodel.UserKey {
	return usermodel.NewUserKey(r.CreatedBy)
}

func (r *Resource) IsOffer() bool {
	return r.Type == Offer
}

func (r *Resource) IsRequest() bool {
	return r.Type == Request
}

func (r *Resource) IsService() bool {
	return r.SubType == ServiceResource
}

func (r *Resource) IsObject() bool {
	return r.SubType == ObjectResource
}