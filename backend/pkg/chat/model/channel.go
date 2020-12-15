package model

import "time"

type Channel struct {
	Key       ChannelKey
	Title     string
	Type      ChannelType
	CreatedAt time.Time
	DeletedAt *time.Time
}

func (c *Channel) GetKey() ChannelKey {
	return c.Key
}
