package readmodel

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/domain"
	"time"
)

type ResourceReadModelBase struct {
	ResourceKey       keys.ResourceKey `gorm:"type:varchar(128);primaryKey;not null"`
	CreatedBy         string           `gorm:"not null;type:varchar(128)"`
	CreatedByVersion  int              `gorm:"not null"`
	CreatedByName     string           `gorm:"not null;type:varchar(128)"`
	CreatedAt         time.Time        `gorm:"not null"`
	UpdatedBy         string           `gorm:"not null;type:varchar(128)"`
	UpdatedByVersion  int              `gorm:"not null"`
	UpdatedByName     string           `gorm:"not null;type:varchar(128)"`
	UpdatedAt         time.Time        `gorm:"not null"`
	GroupSharingCount int              `gorm:"not null"`
	Version           int              `gorm:"not null"`
	Owner             keys.Target      `json:"owner" gorm:"embedded;embeddedPrefix:owner_"`
}

type DbResourceReadModel struct {
	ResourceReadModelBase
	domain.ResourceInfoBase
	domain.ResourceValueEstimation
}

func (d DbResourceReadModel) TableName() string {
	return "resource_read_models"
}

type ResourceReadModel struct {
	ResourceReadModelBase
	domain.ResourceInfo `json:"info"`
}

type ResourceSharingReadModel struct {
	ResourceKey  keys.ResourceKey `gorm:"primaryKey"`
	GroupKey     keys.GroupKey    `gorm:"primaryKey"`
	GroupName    string           `gorm:"not null;type:varchar(128)"`
	Version      int              `gorm:"not null"`
	GroupVersion int              `gorm:"not null"`
}

type ResourceUserNameReadModel struct {
	UserKey  string `gorm:"not null;primaryKey;type:varchar(128)"`
	Username string `gorm:"not null;type:varchar(128)"`
	Version  int    `gorm:"not null"`
}

type ResourceGroupNameReadModel struct {
	GroupKey  string `gorm:"not null;primaryKey;type:varchar(128)"`
	GroupName string `gorm:"not null;type:varchar(128)"`
	Version   int    `gorm:"not null"`
}

type ResourceWithSharingsReadModel struct {
	ResourceReadModel
	Sharings []*ResourceSharingReadModel `json:"sharings"`
}
