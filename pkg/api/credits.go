package api

import "time"

type Credits struct {
	ID        string
	GroupID   string
	Group     *Group
	SentTo    *Target `gorm:"embedded;embeddedPrefix:sent_to_"`
	SentBy    *Target `gorm:"embedded;embeddedPrefix:sent_by_"`
	Amount    time.Duration
	CreatedAt time.Time
	Notes     string
}
