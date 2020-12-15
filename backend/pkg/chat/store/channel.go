package store

import (
	"github.com/commonpool/backend/pkg/chat/model"
	"time"
)

type Channel struct {
	ID        string `gorm:"primaryKey"`
	Title     string
	Type      model.ChannelType
	CreatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (channel *Channel) Map() *model.Channel {
	return &model.Channel{
		Key:       model.NewChannelKey(channel.ID),
		Title:     channel.Title,
		Type:      channel.Type,
		CreatedAt: channel.CreatedAt,
		DeletedAt: channel.DeletedAt,
	}
}

func MapChannel(c *model.Channel) *Channel {
	return &Channel{
		ID:        c.Key.ID,
		Title:     c.Title,
		Type:      c.Type,
		CreatedAt: c.CreatedAt,
		DeletedAt: c.DeletedAt,
	}
}
