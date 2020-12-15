package model

import (
	usermodel "github.com/commonpool/backend/pkg/user/model"
	"time"
)

type ChannelSubscription struct {
	ChannelKey          ChannelKey
	UserKey             usermodel.UserKey
	Name                string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
	LastMessageAt       time.Time
	LastTimeRead        time.Time
	LastMessageChars    string
	LastMessageUserKey  usermodel.UserKey
	LastMessageUserName string
}

func (s *ChannelSubscription) GetKey() ChannelSubscriptionKey {
	return NewChannelSubscriptionKey(
		s.ChannelKey,
		s.UserKey,
	)
}

func (s *ChannelSubscription) GetChannelKey() ChannelKey {
	return s.GetKey().ChannelKey
}

func (s *ChannelSubscription) GetUserKey() usermodel.UserKey {
	return s.GetKey().UserKey
}
