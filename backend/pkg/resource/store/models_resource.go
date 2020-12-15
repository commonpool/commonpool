package store

import (
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
	"time"
)

type Resource struct {
	ID               string                `mapstructure:"id"`
	CreatedAt        time.Time             `mapstructure:"createdAt"`
	UpdatedAt        time.Time             `mapstructure:"updatedAt"`
	DeletedAt        *time.Time            `mapstructure:"deletedAt"`
	Summary          string                `mapstructure:"summary"`
	Description      string                `mapstructure:"description"`
	CreatedBy        string                `mapstructure:"createdBy"`
	Type             resourcemodel.Type    `mapstructure:"type"`
	SubType          resourcemodel.SubType `mapstructure:"subType"`
	ValueInHoursFrom int                   `mapstructure:"valueInHoursFrom"`
	ValueInHoursTo   int                   `mapstructure:"valueInHoursTo"`
}

func mapGraphResourceToResource(dbResultItem *Resource) (*resourcemodel.Resource, error) {

	key, err := resourcemodel.ParseResourceKey(dbResultItem.ID)
	if err != nil {
		return nil, err
	}
	return &resourcemodel.Resource{
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
