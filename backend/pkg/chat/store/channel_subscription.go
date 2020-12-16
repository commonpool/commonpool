package store

import (
	"github.com/commonpool/backend/pkg/chat"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
	"time"
)

type ChannelSubscription struct {
	ChannelID           string `gorm:"primaryKey;not null"`
	UserID              string `gorm:"primaryKey;not null"`
	Name                string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time `sql:"index"`
	LastMessageAt       time.Time
	LastTimeRead        time.Time
	LastMessageChars    string
	LastMessageUserId   string
	LastMessageUserName string
}

func (s *ChannelSubscription) Map() *chat.ChannelSubscription {
	return &chat.ChannelSubscription{
		ChannelKey:          chat.NewChannelKey(s.ChannelID),
		UserKey:             usermodel.NewUserKey(s.UserID),
		Name:                s.Name,
		CreatedAt:           s.CreatedAt,
		UpdatedAt:           s.UpdatedAt,
		DeletedAt:           s.DeletedAt,
		LastMessageAt:       s.LastMessageAt,
		LastTimeRead:        s.LastTimeRead,
		LastMessageChars:    s.LastMessageChars,
		LastMessageUserKey:  usermodel.NewUserKey(s.LastMessageUserId),
		LastMessageUserName: s.LastMessageUserName,
	}
}

func MapChannelSubscription(s *chat.ChannelSubscription) *ChannelSubscription {
	return &ChannelSubscription{
		ChannelID:           s.ChannelKey.String(),
		UserID:              s.UserKey.String(),
		Name:                s.Name,
		CreatedAt:           s.CreatedAt,
		UpdatedAt:           s.UpdatedAt,
		DeletedAt:           s.DeletedAt,
		LastMessageAt:       s.LastMessageAt,
		LastTimeRead:        s.LastTimeRead,
		LastMessageChars:    s.LastMessageChars,
		LastMessageUserId:   s.LastMessageUserKey.String(),
		LastMessageUserName: s.LastMessageUserName,
	}
}
