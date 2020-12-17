package store

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource"
	"time"
)

type Resource struct {
	ID               string           `mapstructure:"id"`
	CreatedAt        time.Time        `mapstructure:"createdAt"`
	UpdatedAt        time.Time        `mapstructure:"updatedAt"`
	DeletedAt        *time.Time       `mapstructure:"deletedAt"`
	Summary          string           `mapstructure:"summary"`
	Description      string           `mapstructure:"description"`
	CreatedBy        string           `mapstructure:"createdBy"`
	Type             resource.Type    `mapstructure:"type"`
	SubType          resource.SubType `mapstructure:"subType"`
	ValueInHoursFrom int              `mapstructure:"valueInHoursFrom"`
	ValueInHoursTo   int              `mapstructure:"valueInHoursTo"`
}

func mapGraphResourceToResource(dbResultItem *Resource) (*resource.Resource, error) {

	key, err := keys.ParseResourceKey(dbResultItem.ID)
	if err != nil {
		return nil, err
	}
	return &resource.Resource{
		Key:              key,
		CreatedAt:        dbResultItem.CreatedAt,
		UpdatedAt:        dbResultItem.UpdatedAt,
		DeletedAt:        dbResultItem.DeletedAt,
		Summary:          dbResultItem.Summary,
		Description:      dbResultItem.Description,
		CreatedBy:        dbResultItem.CreatedBy,
		Type:             dbResultItem.Type,
		SubType:          dbResultItem.SubType,
		ValueInHoursFrom: dbResultItem.ValueInHoursFrom,
		ValueInHoursTo:   dbResultItem.ValueInHoursTo,
	}, nil
}
