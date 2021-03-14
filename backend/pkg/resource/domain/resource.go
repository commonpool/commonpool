package domain

import (
	"fmt"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/keys"
	"strings"
)

type Resource struct {
	key          keys.ResourceKey
	version      int
	changes      []eventsource.Event
	name         string
	description  string
	isNew        bool
	isDeleted    bool
	registeredBy keys.UserKey
	info         ResourceInfo
	sharings     []ResourceSharing
	evaluations  map[keys.UserKey][]ValueEstimations
}

func NewResource(key keys.ResourceKey) *Resource {
	return &Resource{
		key:         key,
		version:     0,
		changes:     []eventsource.Event{},
		isNew:       true,
		isDeleted:   false,
		evaluations: map[keys.UserKey][]ValueEstimations{},
	}
}

func NewFromEvents(key keys.ResourceKey, events []eventsource.Event) *Resource {
	r := NewResource(key)
	for _, event := range events {
		r.on(event, false)
	}
	return r
}

func (r *Resource) Register(
	registeredBy keys.UserKey,
	registeredFor *keys.Target,
	resourceInfo ResourceInfo,
	resourceSharings *keys.GroupKeys,
	values ValueEstimations) error {
	if err := r.assertIsNew(); err != nil {
		return err
	}

	switch resourceInfo.CallType {
	case Offer:
	case Request:
	default:
		return exceptions.ErrBadRequestf("invalid resource call type: %s", resourceInfo.CallType)
	}

	switch resourceInfo.ResourceType {
	case ServiceResource:
	case ObjectResource:
	default:
		return exceptions.ErrBadRequestf("invalid resource type: %s", resourceInfo.ResourceType)
	}

	sanitizedInfo, err := r.sanitizeResourceInfo(resourceInfo)
	if err != nil {
		return err
	}

	if registeredFor == nil {
		return exceptions.ErrBadRequest("registeredFor is required")
	}

	r.raise(NewResourceRegistered(registeredBy, *registeredFor, sanitizedInfo))

	if resourceSharings == nil {
		resourceSharings = keys.NewEmptyGroupKeys()
	}

	err = r.ChangeSharings(registeredBy, resourceSharings)
	if err != nil {
		return err
	}

	err = r.EvaluateResource(registeredBy, values)
	if err != nil {
		return err
	}

	return nil
}

func (r *Resource) handleResourceRegistered(e ResourceRegistered) {
	r.isNew = false
	r.registeredBy = e.RegisteredBy
	r.info = e.ResourceInfo
}

func (r *Resource) ChangeInfo(changedBy keys.UserKey, resourceInfo ResourceInfoUpdate) error {
	if err := r.assertIsNotNew(); err != nil {
		return err
	}
	if err := r.assertNotDeleted(); err != nil {
		return err
	}

	newInfo := r.info.
		WithName(resourceInfo.Name).
		WithDescription(resourceInfo.Description)

	sanitizedInfo, err := r.sanitizeResourceInfo(newInfo)
	if err != nil {
		return err
	}

	if r.info == sanitizedInfo {
		return nil
	}

	evt := NewResourceInfoChanged(changedBy, r.info, sanitizedInfo)
	r.raise(evt)
	return nil
}

func (r *Resource) handleResourceInfoChanged(e ResourceInfoChanged) {
	r.info = e.NewResourceInfo
}

func (r *Resource) EvaluateResource(evaluatedBy keys.UserKey, evaluation ValueEstimations) error {
	evaluations, ok := r.evaluations[evaluatedBy]
	if evaluations == nil {
		evaluations = []ValueEstimations{}
	}

	if len(evaluations) > 0 {
		lastEvaluation := evaluations[len(evaluations)-1]
		if lastEvaluation.Equals(evaluation) {
			return nil
		}
	}

	r.raise(NewResourceEvaluated(evaluatedBy, evaluation, evaluations, !ok))

	return nil
}

func (r *Resource) handleResourceEvaluated(e ResourceEvaluated) {
	r.evaluations[e.EvaluatedBy] = append(r.evaluations[e.EvaluatedBy], e.NewEvaluation)
}

func (r *Resource) ChangeSharings(changedBy keys.UserKey, sharedWith *keys.GroupKeys) error {
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

	evt := NewResourceGroupSharingChanged(changedBy, r.sharings, desiredSharings, toAdd, toDelete, r.info)
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

func (r *Resource) sanitizeResourceInfo(resourceInfo ResourceInfo) (ResourceInfo, error) {
	sanitizedInfo := ResourceInfo{
		ResourceInfoBase: ResourceInfoBase{
			Name:         strings.TrimSpace(resourceInfo.Name),
			Description:  strings.TrimSpace(resourceInfo.Description),
			CallType:     resourceInfo.CallType,
			ResourceType: resourceInfo.ResourceType,
		},
	}

	if sanitizedInfo.Name == "" {
		return ResourceInfo{}, exceptions.ErrValidation("name is required")
	}
	if len(sanitizedInfo.Name) > 64 {
		return ResourceInfo{}, exceptions.ErrValidation("name is too long")
	}
	if sanitizedInfo.Description == "" {
		return ResourceInfo{}, exceptions.ErrValidation("description is required")
	}
	if len(sanitizedInfo.Description) > 2048 {
		return ResourceInfo{}, exceptions.ErrValidation("description is too long")
	}
	return sanitizedInfo, nil
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

func (r *Resource) StreamKey() keys.StreamKey {
	return r.key.StreamKey()
}

func (r *Resource) GetCallType() CallType {
	return r.info.CallType
}

func (r *Resource) GetResourceType() ResourceType {
	return r.info.ResourceType
}

func (r *Resource) GetVersion() int {
	return r.version
}

func (r *Resource) GetChanges() []eventsource.Event {
	return r.changes
}

func (r *Resource) GetKey() keys.ResourceKey {
	return r.key
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
	case ResourceEvaluated:
		o.handleResourceEvaluated(e)
	}
	if !isNew {
		o.version++
	}
}
