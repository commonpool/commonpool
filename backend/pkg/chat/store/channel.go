package store

import (
	"github.com/commonpool/backend/pkg/chat"
	"time"
)

type Channel struct {
	ID        string `gorm:"primaryKey"`
	Title     string
	Type      chat.ChannelType
	CreatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (channel *Channel) Map() *chat.Channel {
	return &chat.Channel{
		Key:       chat.NewChannelKey(channel.ID),
		Title:     channel.Title,
		Type:      channel.Type,
		CreatedAt: channel.CreatedAt,
		DeletedAt: channel.DeletedAt,
	}
}

func MapChannel(c *chat.Channel) *Channel {
	return &Channel{
		ID:        c.Key.ID,
		Title:     c.Title,
		Type:      c.Type,
		CreatedAt: c.CreatedAt,
		DeletedAt: c.DeletedAt,
	}
}
