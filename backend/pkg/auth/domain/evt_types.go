package domain

import (
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/pkg/eventsource"
)

const (
	UserDiscoveredEvent  = "user_discovered"
	UserInfoChangedEvent = "user_info_changed"
)

func RegisterEvents(mapper *eventsource.EventMapper) error {
	for _, evt := range []string{
		UserDiscoveredEvent,
		UserInfoChangedEvent,
	} {
		if err := mapper.RegisterMapper(evt, MapEvent); err != nil {
			return err
		}
	}
	return nil
}

func MapEvent(eventType string, bytes []byte) (eventsource.Event, error) {
	switch eventType {
	case UserDiscoveredEvent:
		dest := UserDiscovered{}
		err := json.Unmarshal(bytes, &dest)
		return dest, err
	case UserInfoChangedEvent:
		dest := UserInfoChanged{}
		err := json.Unmarshal(bytes, &dest)
		return dest, err
	default:
		return nil, fmt.Errorf("invalid command type")
	}
}
