package store

import (
	"github.com/commonpool/backend/pkg/chat/model"
	"github.com/satori/go.uuid"
	"time"
)

type Message struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key"`
	ChannelID      string
	MessageType    model.MessageType
	MessageSubType model.MessageSubType
	SentById       string
	SentByUsername string
	SentAt         time.Time
	Text           string
	Blocks         string `gorm:"type:jsonb"`
	Attachments    string `gorm:"type:jsonb"`
	VisibleToUser  *string
}
