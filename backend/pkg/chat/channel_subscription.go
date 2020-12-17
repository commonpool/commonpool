package chat

import (
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type ChannelSubscription struct {
	ChannelKey          keys.ChannelKey
	UserKey             keys.UserKey
	Name                string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
	LastMessageAt       time.Time
	LastTimeRead        time.Time
	LastMessageChars    string
	LastMessageUserKey  keys.UserKey
	LastMessageUserName string
}

func (s *ChannelSubscription) GetKey() keys.ChannelSubscriptionKey {
	return keys.NewChannelSubscriptionKey(
		s.ChannelKey,
		s.UserKey,
	)
}

func (s *ChannelSubscription) GetChannelKey() keys.ChannelKey {
	return s.GetKey().ChannelKey
}

func (s *ChannelSubscription) GetUserKey() keys.UserKey {
	return s.GetKey().UserKey
}
