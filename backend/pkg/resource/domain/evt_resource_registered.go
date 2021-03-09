package domain

import (
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
)

type ResourceRegisteredPayload struct {
	RegisteredBy  keys.UserKey `json:"registered_by"`
	RegisteredFor keys.Target  `json:"registered_for"`
	ResourceInfo  ResourceInfo `json:"resource_info"`
}

type ResourceRegistered struct {
	eventsource.EventEnvelope
	ResourceRegisteredPayload `json:"payload"`
}

func NewResourceRegistered(registeredBy keys.UserKey, registeredFor keys.Target, resourceInfo ResourceInfo) ResourceRegistered {
	return ResourceRegistered{
		eventsource.NewEventEnvelope(ResourceRegisteredEvent, 1),
		ResourceRegisteredPayload{
			RegisteredBy:  registeredBy,
			RegisteredFor: registeredFor,
			ResourceInfo:  resourceInfo,
		},
	}
}
