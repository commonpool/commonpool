package chat

import (
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type Channel struct {
	Key       keys.ChannelKey
	Title     string
	Type      ChannelType
	CreatedAt time.Time
	DeletedAt *time.Time
}

func (c *Channel) GetKey() keys.ChannelKey {
	return c.Key
}
