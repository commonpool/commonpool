package readmodel

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/domain"
	"time"
)

type ResourceReadModelBase struct {
	ResourceKey       keys.ResourceKey `gorm:"type:varchar(128);primaryKey;not null" json:"resource_key"`
	CreatedBy         string           `gorm:"not null;type:varchar(128)" json:"created_by"`
	CreatedByVersion  int              `gorm:"not null" json:"created_by_version"`
	CreatedByName     string           `gorm:"not null;type:varchar(128)" json:"created_by_name"`
	CreatedAt         time.Time        `gorm:"not null" json:"created_at"`
	UpdatedBy         string           `gorm:"not null;type:varchar(128)" json:"updated_by"`
	UpdatedByVersion  int              `gorm:"not null" json:"updated_by_version"`
	UpdatedByName     string           `gorm:"not null;type:varchar(128)" json:"updated_by_name"`
	UpdatedAt         time.Time        `gorm:"not null" json:"updated_at"`
	GroupSharingCount int              `gorm:"not null" json:"group_sharing_count"`
	Version           int              `gorm:"not null" json:"version"`
	Owner             keys.Target      `json:"owner" gorm:"embedded;embeddedPrefix:owner_" json:"owner"`
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
	ResourceKey  keys.ResourceKey `gorm:"primaryKey" json:"resource_key"`
	GroupKey     keys.GroupKey    `gorm:"primaryKey" json:"group_key"`
	GroupName    string           `gorm:"not null;type:varchar(128)" json:"group_name"`
	Version      int              `gorm:"not null" json:"version"`
	GroupVersion int              `gorm:"not null" json:"group_version"`
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
