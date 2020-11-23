package model

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Message struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key"`
	TopicID        string    `gorm:"primary_key"`
	MessageType    MessageType
	MessageSubType MessageSubType
	UserID         string
	BotID          string
	SentAt         time.Time
	Text           string
	Blocks         string `gorm:"type:jsonb"`
	Attachments    string `gorm:"type:jsonb"`
	IsPersonal     bool
	SentBy         string
	SentByUsername string
}

func (m *Message) GetAuthorKey() UserKey {
	return NewUserKey(m.SentBy)
}
