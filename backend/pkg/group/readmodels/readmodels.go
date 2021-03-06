package readmodels

import "time"

type GroupReadModel struct {
	Version     int       `gorm:"not null"`
	GroupKey    string    `gorm:"not null;type:varchar(128);primaryKey"`
	Name        string    `gorm:"not null;type:varchar(128)"`
	Description string    `gorm:"not null;type:varchar(2048)"`
	CreatedBy   string    `gorm:"not null;type:varchar(128)"`
	CreatedAt   time.Time `gorm:"not null"`
}

type MembershipReadModel struct {
	Version          int     `gorm:"not null"`
	GroupKey         string  `gorm:"not null;type:varchar(128);primaryKey"`
	GroupName        string  `gorm:"type:varchar(128)"`
	UserKey          string  `gorm:"not null;type:varchar(128);primaryKey"`
	IsOwner          bool    `gorm:"not null"`
	IsAdmin          bool    `gorm:"not null"`
	IsMember         bool    `gorm:"not null"`
	GroupConfirmed   bool    `gorm:"not null"`
	GroupConfirmedBy *string `gorm:"type:varchar(128)"`
	GroupConfirmedAt *time.Time
	UserConfirmed    bool `gorm:"not null"`
	UserConfirmedAt  *time.Time
}
