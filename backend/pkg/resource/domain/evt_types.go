package domain

import (
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/pkg/eventsource"
)

const (
	ResourceRegisteredEvent          = "resource_registered"
	ResourceGroupSharingChangedEvent = "resource_group_sharing_changed"
	ResourceInfoChangedEvent         = "resource_info_changed"
	ResourceDeletedEvent             = "resource_deleted"
)

var AllEventTypes = []string{
	ResourceRegisteredEvent,
	ResourceGroupSharingChangedEvent,
	ResourceInfoChangedEvent,
	ResourceDeletedEvent,
}

func RegisterEvents(eventMapper *eventsource.EventMapper) error {
	for _, eventType := range AllEventTypes {
		if err := eventMapper.RegisterMapper(eventType, MapEvent); err != nil {
			return err
		}
	}
	return nil
}

func MapEvent(eventType string, bytes []byte) (eventsource.Event, error) {
	switch eventType {
	case ResourceRegisteredEvent:
		var dest ResourceRegistered
		err := json.Unmarshal(bytes, &dest)
		return dest, err
	case ResourceGroupSharingChangedEvent:
		var dest ResourceGroupSharingChanged
		err := json.Unmarshal(bytes, &dest)
		return dest, err
	case ResourceInfoChangedEvent:
		var dest ResourceInfoChanged
		err := json.Unmarshal(bytes, &dest)
		return dest, err
	case ResourceDeletedEvent:
		var dest ResourceDeleted
		err := json.Unmarshal(bytes, &dest)
		return dest, err
	default:
		return nil, fmt.Errorf("unexpected event type")
	}
}
