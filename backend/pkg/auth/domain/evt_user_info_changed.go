package domain

import (
	"github.com/commonpool/backend/pkg/eventsource"
	"time"
)

type UserInfoChangedPayload struct {
	OldUserInfo UserInfo `json:"old_user_info"`
	NewUserInfo UserInfo `json:"new_user_info"`
}

// UserInfoChanged Emitted when a user changes his profile information
type UserInfoChanged struct {
	eventsource.EventEnvelope
	UserInfoChangedPayload `json:"payload"`
}

var _ eventsource.Event = &UserInfoChanged{}

func NewUserInfoChanged(oldUserInfo, newUserInfo UserInfo) UserInfoChanged {
	return UserInfoChanged{
		EventEnvelope: eventsource.EventEnvelope{
			EventTime: time.Now(),
			EventType: UserInfoChangedEvent,
		},
		UserInfoChangedPayload: UserInfoChangedPayload{
			OldUserInfo: oldUserInfo,
			NewUserInfo: newUserInfo,
		},
	}
}
