package domain

import (
	"github.com/commonpool/backend/pkg/eventsource"
	"time"
)

type UserDiscoveredPayload struct {
	UserInfo UserInfo `json:"user_info"`
}

// UserDiscovered Emitted when a user first logs in
// to the system
type UserDiscovered struct {
	eventsource.EventEnvelope
	UserDiscoveredPayload `json:"payload"`
}

var _ eventsource.Event = &UserDiscovered{}

func NewUserDiscovered(userInfo UserInfo) UserDiscovered {
	return UserDiscovered{
		EventEnvelope: eventsource.EventEnvelope{
			EventTime:     time.Now().UTC(),
			EventType:     UserDiscoveredEvent,
			CorrelationID: "",
			EventID:       "",
			AggregateID:   "",
			AggregateType: "",
			EventVersion:  0,
			SequenceNo:    0,
		},
		UserDiscoveredPayload: UserDiscoveredPayload{
			UserInfo: userInfo,
		},
	}
}
