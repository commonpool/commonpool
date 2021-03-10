package readmodel

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/domain"
	"time"
)

type ResourceReadModelBase struct {
	ResourceKey       keys.ResourceKey `gorm:"type:varchar(128);primaryKey;not null" json:"resourceId"`
	CreatedBy         string           `gorm:"not null;type:varchar(128)" json:"createdBy"`
	CreatedByVersion  int              `gorm:"not null" json:"createdByVersion"`
	CreatedByName     string           `gorm:"not null;type:varchar(128)" json:"createdByName"`
	CreatedAt         time.Time        `gorm:"not null" json:"createdAt"`
	UpdatedBy         string           `gorm:"not null;type:varchar(128)" json:"updatedBy"`
	UpdatedByVersion  int              `gorm:"not null" json:"updatedByVersion"`
	UpdatedByName     string           `gorm:"not null;type:varchar(128)" json:"updatedByName"`
	UpdatedAt         time.Time        `gorm:"not null" json:"updatedAt"`
	GroupSharingCount int              `gorm:"not null" json:"groupSharingCount"`
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
	ResourceKey  keys.ResourceKey `gorm:"primaryKey" json:"resourceId"`
	GroupKey     keys.GroupKey    `gorm:"primaryKey" json:"groupId"`
	GroupName    string           `gorm:"not null;type:varchar(128)" json:"groupName"`
	Version      int              `gorm:"not null" json:"version"`
	GroupVersion int              `gorm:"not null" json:"groupVersion"`
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
