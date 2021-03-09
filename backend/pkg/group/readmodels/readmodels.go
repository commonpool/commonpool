package readmodels

import (
	"github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type GroupReadModel struct {
	Version     int           `gorm:"not null"`
	GroupKey    keys.GroupKey `gorm:"not null;type:varchar(128);primaryKey"`
	Name        string        `gorm:"not null;type:varchar(128)"`
	Description string        `gorm:"not null;type:varchar(2048)"`
	CreatedBy   string        `gorm:"not null;type:varchar(128)"`
	CreatedAt   time.Time     `gorm:"not null"`
}

type MembershipReadModel struct {
	Version          int           `gorm:"not null"`
	GroupKey         keys.GroupKey `gorm:"not null;type:varchar(128);primaryKey"`
	GroupName        string        `gorm:"type:varchar(128)"`
	UserKey          keys.UserKey  `gorm:"not null;type:varchar(128);primaryKey"`
	IsOwner          bool          `gorm:"not null"`
	IsAdmin          bool          `gorm:"not null"`
	IsMember         bool          `gorm:"not null"`
	GroupConfirmed   bool          `gorm:"not null"`
	GroupConfirmedBy *string       `gorm:"type:varchar(128)"`
	GroupConfirmedAt *time.Time
	UserConfirmed    bool `gorm:"not null"`
	UserConfirmedAt  *time.Time
	Status           domain.MembershipStatus `gorm:"not null"`
	UserVersion      int                     `gorm:"not null"`
	UserName         string                  `gorm:"type:varchar(128)"`
}

type DBGroupUserReadModel struct {
	UserKey keys.UserKey `gorm:"primaryKey"`
	Name    string
	Version int
}

func (d DBGroupUserReadModel) TableName() string {
	return "group_users"
}
