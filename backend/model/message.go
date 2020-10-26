package model

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Message struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key"`
	TopicId  uuid.UUID `gorm:"type:uuid"`
	UserID   string
	AuthorID string
	SentAt   time.Time
	Content  string
}

func (m *Message) GetAuthorKey() UserKey {
	return NewUserKey(m.AuthorID)
}
