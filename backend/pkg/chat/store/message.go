package store

import (
	"github.com/commonpool/backend/pkg/chat/chatmodel"
	"github.com/satori/go.uuid"
	"time"
)

type Message struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key"`
	ChannelID      string
	MessageType    chatmodel.MessageType
	MessageSubType chatmodel.MessageSubType
	SentById       string
	SentByUsername string
	SentAt         time.Time
	Text           string
	Blocks         string `gorm:"type:jsonb"`
	Attachments    string `gorm:"type:jsonb"`
	VisibleToUser  *string
}
