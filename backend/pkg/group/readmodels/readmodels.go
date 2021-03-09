package readmodels

import (
	"github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type GroupReadModel struct {
	Version     int           `gorm:"not null" json:"version"`
	GroupKey    keys.GroupKey `gorm:"not null;type:varchar(128);primaryKey" json:"group_key"`
	Name        string        `gorm:"not null;type:varchar(128)" json:"name"`
	Description string        `gorm:"not null;type:varchar(2048)" json:"description"`
	CreatedBy   string        `gorm:"not null;type:varchar(128)" json:"created_by"`
	CreatedAt   time.Time     `gorm:"not null" json:"created_at"`
}

type MembershipReadModel struct {
	Version          int                     `gorm:"not null" json:"version"`
	GroupKey         keys.GroupKey           `gorm:"not null;type:varchar(128);primaryKey" json:"group_key"`
	GroupName        string                  `gorm:"type:varchar(128)" json:"group_name"`
	UserKey          keys.UserKey            `gorm:"not null;type:varchar(128);primaryKey" json:"user_key"`
	IsOwner          bool                    `gorm:"not null" json:"is_owner"`
	IsAdmin          bool                    `gorm:"not null" json:"is_admin"`
	IsMember         bool                    `gorm:"not null" json:"is_member"`
	GroupConfirmed   bool                    `gorm:"not null" json:"group_confirmed"`
	GroupConfirmedBy *string                 `gorm:"type:varchar(128)" json:"group_confirmed_by"`
	GroupConfirmedAt *time.Time              `json:"group_confirmed_at"`
	UserConfirmed    bool                    `gorm:"not null" json:"user_confirmed"`
	UserConfirmedAt  *time.Time              `json:"user_confirmed_at"`
	Status           domain.MembershipStatus `gorm:"not null" json:"status"`
	UserVersion      int                     `gorm:"not null" json:"user_version"`
	UserName         string                  `gorm:"type:varchar(128)" json:"user_name"`
}

type DBGroupUserReadModel struct {
	UserKey keys.UserKey `gorm:"primaryKey"`
	Name    string
	Version int
}

func (d DBGroupUserReadModel) TableName() string {
	return "group_users"
}
