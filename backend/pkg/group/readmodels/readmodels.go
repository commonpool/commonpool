package readmodels

import (
	"github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type GroupReadModel struct {
	Version     int           `gorm:"not null" json:"version"`
	GroupKey    keys.GroupKey `gorm:"not null;type:varchar(128);primaryKey" json:"groupId"`
	Name        string        `gorm:"not null;type:varchar(128)" json:"name"`
	Description string        `gorm:"not null;type:varchar(2048)" json:"description"`
	CreatedBy   string        `gorm:"not null;type:varchar(128)" json:"createdBy"`
	CreatedAt   time.Time     `gorm:"not null" json:"createdAt"`
}

func (g GroupReadModel) Target() *keys.Target {
	return g.GroupKey.Target()
}

type MembershipReadModel struct {
	Version          int                     `gorm:"not null" json:"version"`
	GroupKey         keys.GroupKey           `gorm:"not null;type:varchar(128);primaryKey" json:"groupId"`
	GroupName        string                  `gorm:"type:varchar(128)" json:"groupName"`
	UserKey          keys.UserKey            `gorm:"not null;type:varchar(128);primaryKey" json:"userId"`
	IsOwner          bool                    `gorm:"not null" json:"isOwner"`
	IsAdmin          bool                    `gorm:"not null" json:"isAdmin"`
	IsMember         bool                    `gorm:"not null" json:"isMember"`
	GroupConfirmed   bool                    `gorm:"not null" json:"groupConfirmed"`
	GroupConfirmedBy *string                 `gorm:"type:varchar(128)" json:"groupConfirmedBy"`
	GroupConfirmedAt *time.Time              `json:"groupConfirmedAt"`
	UserConfirmed    bool                    `gorm:"not null" json:"userConfirmed"`
	UserConfirmedAt  *time.Time              `json:"userConfirmedAt"`
	Status           domain.MembershipStatus `gorm:"not null" json:"status"`
	UserVersion      int                     `gorm:"not null" json:"userVersion"`
	UserName         string                  `gorm:"type:varchar(128)" json:"userName"`
	CreatedBy        string                  `gorm:"type:varchar(128)" json:"createdBy"`
	CreatedByName    string                  `gorm:"type:varchar(128)" json:"createdByName"`
	CreatedByVersion int                     `gorm:"not null" json:"createdByVersion"`
	CreatedAt        time.Time               `json:"createdAt"`
}

type DBGroupUserReadModel struct {
	UserKey keys.UserKey `gorm:"primaryKey"`
	Name    string
	Version int
}

func (d DBGroupUserReadModel) TableName() string {
	return "group_users"
}
