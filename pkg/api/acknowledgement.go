package api

import "time"

type AcknowledgementType string

const (
	ThanksObjectGift  = "thanks-gift-object"
	ThanksServiceGift = "thanks-gift-service"
	ThanksObjectLent  = "thanks-lent-object"
	Other  = "other"
)

type Acknowledgement struct {
	ID        string
	GroupID   string
	Group     *Group
	SentTo    *Target `gorm:"embedded;embeddedPrefix:sent_to_"`
	SentBy    *Target `gorm:"embedded;embeddedPrefix:sent_by_"`
	CreatedAt time.Time
	Type      AcknowledgementType
	Notes     string
}
