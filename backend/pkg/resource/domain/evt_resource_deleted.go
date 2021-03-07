package domain

import (
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
)

type ResourceDeletedPayload struct {
	DeletedBy keys.UserKey `json:"registered_by"`
}

type ResourceDeleted struct {
	eventsource.EventEnvelope
	ResourceDeletedPayload `json:"payload"`
}

func NewResourceDeleted(deletedBy keys.UserKey) ResourceDeleted {
	return ResourceDeleted{
		eventsource.NewEventEnvelope(ResourceDeletedEvent, 1),
		ResourceDeletedPayload{
			DeletedBy: deletedBy,
		},
	}
}
