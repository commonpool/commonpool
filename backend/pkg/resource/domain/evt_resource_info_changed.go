package domain

import (
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
)

type ResourceInfoChangedPayload struct {
	ChangedBy       keys.UserKey `json:"registered_by"`
	OldResourceInfo ResourceInfo `json:"old_resource_info"`
	NewResourceInfo ResourceInfo `json:"new_resource_info"`
}

type ResourceInfoChanged struct {
	eventsource.EventEnvelope
	ResourceInfoChangedPayload `json:"payload"`
}

func NewResourceInfoChanged(changedBy keys.UserKey, oldResourceInfo, newResourceInfo ResourceInfo) ResourceInfoChanged {
	return ResourceInfoChanged{
		eventsource.NewEventEnvelope(ResourceInfoChangedEvent, 1),
		ResourceInfoChangedPayload{
			ChangedBy:       changedBy,
			OldResourceInfo: oldResourceInfo,
			NewResourceInfo: newResourceInfo,
		},
	}
}
