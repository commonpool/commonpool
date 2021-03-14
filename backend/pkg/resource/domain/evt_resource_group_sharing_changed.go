package domain

import (
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
)

type ResourceSharings []ResourceSharing

type ResourceGroupSharingChangedPayload struct {
	ChangedBy           keys.UserKey     `json:"registeredBy"`
	OldResourceSharings ResourceSharings `json:"oldResourceSharings"`
	NewResourceSharings ResourceSharings `json:"newResourceSharings"`
	AddedSharings       ResourceSharings `json:"addedSharings"`
	RemovedSharings     ResourceSharings `json:"removedSharings"`
	ResourceInfo        ResourceInfo     `json:"resourceInfo"`
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
	removedSharings ResourceSharings,
	resourceInfo ResourceInfo) ResourceGroupSharingChanged {
	return ResourceGroupSharingChanged{
		eventsource.NewEventEnvelope(ResourceGroupSharingChangedEvent, 1),
		ResourceGroupSharingChangedPayload{
			ChangedBy:           changedBy,
			OldResourceSharings: oldResourceSharings,
			NewResourceSharings: newResourceSharings,
			AddedSharings:       addedSharings,
			RemovedSharings:     removedSharings,
			ResourceInfo:        resourceInfo,
		},
	}
}
