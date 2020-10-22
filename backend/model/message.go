package model

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Message struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key"`
	CreatedAt  time.Time
	SenderId   string
	ReceiverId string
	ThreadId   string
	Content    string
}
