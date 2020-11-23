package model

import (
	"time"
)

type Thread struct {
	UserID              string `gorm:"primaryKey"`
	TopicID             string `gorm:"primaryKey"`
	CreatedAt           time.Time
	LastMessageAt       time.Time
	LastTimeRead        time.Time
	LastMessageChars    string
	LastMessageUserId   string
	LastMessageUserName string
}

func (t *Thread) GetKey() ThreadKey {
	return NewThreadKey(NewTopicKey(t.TopicID), NewUserKey(t.UserID))
}
