package domain

import (
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
)

type ResourceSharings []ResourceSharing

type ResourceGroupSharingChangedPayload struct {
	ChangedBy           keys.UserKey     `json:"registered_by"`
	OldResourceSharings ResourceSharings `json:"old_resource_sharings"`
	NewResourceSharings ResourceSharings `json:"new_resource_sharings"`
	AddedSharings       ResourceSharings `json:"added_sharings"`
	RemovedSharings     ResourceSharings `json:"removed_sharings"`
}

type ResourceGroupSharingChanged struct {
	eventsource.EventEnvelope
	ResourceGroupSharingChangedPayload `json:"payload"`
}

func NewResourceGroupSharingChanged(
	changedBy keys.UserKey,
	oldResourceSharings,
	newResourceSharings,
	addedSharings,
	removedSharings ResourceSharings) ResourceGroupSharingChanged {
	return ResourceGroupSharingChanged{
		eventsource.NewEventEnvelope(ResourceGroupSharingChangedEvent, 1),
		ResourceGroupSharingChangedPayload{
			ChangedBy:           changedBy,
			OldResourceSharings: oldResourceSharings,
			NewResourceSharings: newResourceSharings,
			AddedSharings:       addedSharings,
			RemovedSharings:     removedSharings,
		},
	}
}
