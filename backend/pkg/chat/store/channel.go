package store

import (
	"github.com/commonpool/backend/pkg/chat/chatmodel"
	"time"
)

type Channel struct {
	ID        string `gorm:"primaryKey"`
	Title     string
	Type      chatmodel.ChannelType
	CreatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (channel *Channel) Map() *chatmodel.Channel {
	return &chatmodel.Channel{
		Key:       chatmodel.NewChannelKey(channel.ID),
		Title:     channel.Title,
		Type:      channel.Type,
		CreatedAt: channel.CreatedAt,
		DeletedAt: channel.DeletedAt,
	}
}

func MapChannel(c *chatmodel.Channel) *Channel {
	return &Channel{
		ID:        c.Key.ID,
		Title:     c.Title,
		Type:      c.Type,
		CreatedAt: c.CreatedAt,
		DeletedAt: c.DeletedAt,
	}
}
