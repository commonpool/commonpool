package readmodel

import (
	"github.com/commonpool/backend/pkg/resource/domain"
	"time"
)

type ResourceReadModel struct {
	ResourceKey      string `gorm:"type:varchar(128);primaryKey;not null"`
	ResourceName     string `gorm:"not null;type:varchar(128)"`
	Description      string `gorm:"not null"`
	CreatedBy        string `gorm:"not null;type:varchar(128)"`
	CreatedByVersion int
	CreatedByName    string    `gorm:"not null;type:varchar(128)"`
	CreatedAt        time.Time `gorm:"not null"`
	UpdatedBy        string    `gorm:"not null;type:varchar(128)"`
	UpdatedByVersion int
	UpdatedByName    string    `gorm:"not null;type:varchar(128)"`
	UpdatedAt        time.Time `gorm:"not null"`
	domain.ResourceValueEstimation
	GroupSharingCount int `gorm:"not null"`
	Version           int `gorm:"not null"`
}

type ResourceSharingReadModel struct {
	ResourceKey  string `gorm:"primaryKey"`
	GroupKey     string `gorm:"primaryKey"`
	GroupName    string `gorm:"not null;type:varchar(128)"`
	Version      int    `gorm:"not null"`
	GroupVersion int
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
