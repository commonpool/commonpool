package domain

import (
	"fmt"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
)

type Resource struct {
	key          keys.ResourceKey
	version      int
	changes      []eventsource.Event
	name         string
	description  string
	resourceType ResourceType
	isNew        bool
	isDeleted    bool
	registeredBy keys.UserKey
	info         ResourceInfo
	sharings     []ResourceSharing
}

func NewResource(key keys.ResourceKey) *Resource {
	return &Resource{
		key:       key,
		version:   0,
		changes:   []eventsource.Event{},
		isNew:     true,
		isDeleted: false,
	}
}

func NewFromEvents(key keys.ResourceKey, events []eventsource.Event) *Resource {
	r := NewResource(key)
	for _, event := range events {
		r.on(event, false)
	}
	return r
}

func (r *Resource) Register(registeredBy keys.UserKey, registeredFor domain.Target, resourceType ResourceType, resourceInfo ResourceInfo, resourceSharings keys.GroupKeys) error {
	if err := r.assertIsNew(); err != nil {
		return err
	}
	evt := NewResourceRegistered(registeredBy, registeredFor, resourceType, resourceInfo)
	r.raise(evt)

	return r.ChangeSharings(registeredBy, resourceSharings)
}

func (r *Resource) handleResourceRegistered(e ResourceRegistered) {
	r.isNew = false
	r.registeredBy = e.RegisteredBy
	r.resourceType = e.ResourceType
	r.info = e.ResourceInfo
}

func (r *Resource) ChangeInfo(changedBy keys.UserKey, resourceInfo ResourceInfo) error {
	if err := r.assertIsNotNew(); err != nil {
		return err
	}
	if err := r.assertNotDeleted(); err != nil {
		return err
	}
	if r.info == resourceInfo {
		return nil
	}
	evt := NewResourceInfoChanged(changedBy, r.info, resourceInfo)
	r.raise(evt)
	return nil
}

func (r *Resource) handleResourceInfoChanged(e ResourceInfoChanged) {
	r.info = e.NewResourceInfo
}

func (r *Resource) ChangeSharings(changedBy keys.UserKey, sharedWith keys.GroupKeys) error {
	if err := r.assertIsNotNew(); err != nil {
		return err
	}
	if err := r.assertNotDeleted(); err != nil {
		return err
	}

	desiredVisitedMap := map[keys.GroupKey]bool{}
	var desiredSharings []ResourceSharing
	var toAdd []ResourceSharing
	for _, desiredSharing := range sharedWith.Items {
		if desiredVisitedMap[desiredSharing] {
			continue
		}
		desiredVisitedMap[desiredSharing] = true
		desiredSharings = append(desiredSharings, ResourceSharing{
			GroupKey: desiredSharing,
		})
		found := false
		for _, sharing := range r.sharings {
			if sharing.GroupKey == desiredSharing {
				found = true
				break
			}
		}
		if !found {
			toAdd = append(toAdd, ResourceSharing{GroupKey: desiredSharing})
		}
	}

	var toDelete []ResourceSharing
	for _, currentSharing := range r.sharings {
		found := false
		for _, desiredItem := range sharedWith.Items {
			if desiredItem == currentSharing.GroupKey {
				found = true
				break
			}
		}
		if !found {
			toDelete = append(toDelete, currentSharing)
		}
	}

	if len(toAdd) == 0 && len(toDelete) == 0 {
		return nil
	}

	evt := NewResourceGroupSharingChanged(changedBy, r.sharings, desiredSharings, toAdd, toDelete)
	r.raise(evt)
	return nil
}

func (r *Resource) handleResourceGroupSharingChanged(e ResourceGroupSharingChanged) {
	r.sharings = e.NewResourceSharings
}

func (r *Resource) Delete(deletedBy keys.UserKey) error {
	if err := r.assertIsNotNew(); err != nil {
		return err
	}
	if r.isDeleted {
		return nil
	}
	evt := NewResourceDeleted(deletedBy)
	r.raise(evt)
	return nil
}

func (r *Resource) handleResourceDeleted(e ResourceDeleted) {
	r.isDeleted = true
}

func (r *Resource) assertIsNew() error {
	if !r.isNew {
		return fmt.Errorf("resource is not new")
	}
	return nil
}

func (r *Resource) assertNotDeleted() error {
	if r.isDeleted {
		return fmt.Errorf("resource is deleted")
	}
	return nil
}

func (r *Resource) assertIsNotNew() error {
	if r.isNew {
		return fmt.Errorf("resource is not new")
	}
	return nil
}

func (r *Resource) GetVersion() int {
	return r.version
}

func (r *Resource) GetChanges() []eventsource.Event {
	return r.changes
}

func (r *Resource) MarkAsCommitted() {
	r.version += len(r.changes)
	r.changes = []eventsource.Event{}
}

func (o *Resource) raise(event eventsource.Event) {
	o.changes = append(o.changes, event)
	o.on(event, true)
}

func (o *Resource) on(evt eventsource.Event, isNew bool) {
	switch e := evt.(type) {
	case ResourceRegistered:
		o.handleResourceRegistered(e)
	case ResourceInfoChanged:
		o.handleResourceInfoChanged(e)
	case ResourceGroupSharingChanged:
		o.handleResourceGroupSharingChanged(e)
	case ResourceDeleted:
		o.handleResourceDeleted(e)
	}
	if !isNew {
		o.version++
	}
}
