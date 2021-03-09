package domain

import (
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/keys"
	"strings"
)

type Group struct {
	aggregateType string
	key           keys.GroupKey
	changes       []eventsource.Event
	version       int
	isNew         bool
	info          GroupInfo
	memberships   *Memberships
}

func (g Group) GetGroupKey() keys.GroupKey {
	return g.key
}

func NewGroup(key keys.GroupKey) *Group {
	return &Group{
		aggregateType: "group",
		key:           key,
		changes:       []eventsource.Event{},
		version:       0,
		isNew:         true,
		info:          GroupInfo{},
		memberships:   NewMemberships([]*Membership{}),
	}
}

func NewFromEvents(key keys.GroupKey, events []eventsource.Event) *Group {
	group := NewGroup(key)
	group.key = key
	for _, event := range events {
		group.on(event, false)
	}
	return group
}

func (o *Group) CreateGroup(createdBy keys.UserKey, groupInfo GroupInfo) error {
	if err := o.assertIsNew(); err != nil {
		return fmt.Errorf("cannot create group: %v", err)
	}

	groupInfo, err := validateGroupInfo(groupInfo)
	if err != nil {
		return fmt.Errorf("cannot create group: %v", err)
	}
	o.raise(NewGroupCreated(createdBy, groupInfo))

	newStatus := ApprovedMembershipStatus
	o.raise(NewMembershipStatusChanged(
		createdBy,
		createdBy,
		nil,
		&newStatus,
		None,
		Owner,
		true,
		false,
		o.info.Name))

	return nil
}

func (o *Group) handleGroupCreated(e GroupCreated) {
	o.isNew = false
	o.info = e.GroupInfo
}

func (o *Group) ChangeInfo(changedBy keys.UserKey, groupInfo GroupInfo) error {
	if err := o.assertNotNew(); err != nil {
		return fmt.Errorf("cannot change group info: %v", err)
	}

	groupInfo, err := validateGroupInfo(groupInfo)
	if err != nil {
		return fmt.Errorf("cannot create group: %v", err)
	}

	if groupInfo == o.info {
		return nil
	}

	o.raise(NewGroupInfoChanged(changedBy, o.info, groupInfo))

	return nil

}

func (o *Group) handleGroupInfoChanged(e GroupInfoChanged) {
	o.info = e.NewGroupInfo
}

func (o *Group) JoinGroup(requestedBy keys.UserKey, memberKey keys.UserKey) error {
	membershipKey := keys.NewMembershipKey(o, memberKey)
	m, ok := o.memberships.GetMembership(membershipKey)

	if !ok {
		if requestedBy != memberKey {
			if err := o.assertIsAdmin(requestedBy); err != nil {
				return err
			}
			if err := o.assertIsActive(requestedBy); err != nil {
				return err
			}
			newStatus := PendingUserMembershipStatus
			o.raise(NewMembershipStatusChanged(
				requestedBy,
				memberKey,
				nil,
				&newStatus,
				None,
				None,
				true,
				false,
				o.info.Name))
		} else {
			newStatus := PendingGroupMembershipStatus
			o.raise(NewMembershipStatusChanged(
				requestedBy,
				memberKey,
				nil,
				&newStatus,
				None,
				None,
				true,
				false,
				o.info.Name))
		}
	} else {

		if requestedBy != memberKey {
			if err := o.assertIsAdmin(requestedBy); err != nil {
				return err
			}
			if err := o.assertIsActive(requestedBy); err != nil {
				return err
			}
			if m.HasGroupConfirmed() {
				return nil
			}
		} else {
			if m.HasUserConfirmed() {
				return nil
			}
		}
		oldStatus := m.Status
		newStatus := ApprovedMembershipStatus
		o.raise(NewMembershipStatusChanged(
			requestedBy,
			memberKey,
			&oldStatus,
			&newStatus,
			None,
			Member,
			false,
			false,
			o.info.Name))
	}

	return nil

}

func (o *Group) CancelMembership(requestedBy, memberKey keys.UserKey) error {

	membership, ok := o.GetMembership(memberKey)
	if !ok {
		return exceptions.ErrForbidden
	}

	requesterMembership, ok := o.GetMembership(requestedBy)
	if !ok {
		return exceptions.ErrForbidden
	}

	if requestedBy != memberKey {
		if err := o.assertIsAdmin(requestedBy); err != nil {
			return err
		}
		if err := o.assertIsActive(requestedBy); err != nil {
			return err
		}
		if requesterMembership.PermissionLevel < membership.PermissionLevel {
			return exceptions.ErrForbidden
		}
	}

	oldStatus := membership.Status
	o.raise(NewMembershipStatusChanged(
		requestedBy,
		memberKey,
		&oldStatus,
		nil,
		membership.PermissionLevel,
		None,
		false,
		true,
		o.info.Name))

	return nil
}

func (o *Group) handleMembershipStatusChange(e MembershipStatusChanged) {
	membershipKey := keys.NewMembershipKey(o, e.MemberKey)
	if e.IsNewMembership {
		m := &Membership{
			PermissionLevel: e.NewPermissions,
			Key:             membershipKey,
			Status:          *e.NewStatus,
		}
		o.memberships.Items = append(o.memberships.Items, m)
	} else if e.NewStatus == nil {
		o.memberships.RemoveMembership(membershipKey)
	} else {
		m, _ := o.memberships.GetMembership(keys.NewMembershipKey(o, e.MemberKey))
		m.Status = *e.NewStatus
		m.PermissionLevel = e.NewPermissions
	}
}

func (o *Group) AssignPermission(grantedBy keys.UserKey, grantedTo keys.UserKey, permission PermissionLevel) error {

	if err := o.assertIsAdmin(grantedBy); err != nil {
		return err
	}

	if err := o.assertIsActive(grantedBy); err != nil {
		return err
	}

	grantedByMembership, ok := o.GetMembership(grantedBy)
	if !ok {
		return exceptions.ErrMembershipNotFound
	}

	grantedToMembership, ok := o.GetMembership(grantedTo)
	if !ok {
		return exceptions.ErrMembershipNotFound
	}

	if permission > grantedByMembership.PermissionLevel {
		return exceptions.ErrForbidden
	}

	if grantedToMembership.PermissionLevel == permission {
		return nil
	}

	status := grantedToMembership.Status
	o.raise(NewMembershipStatusChanged(
		grantedBy,
		grantedTo,
		&status,
		&status,
		grantedToMembership.PermissionLevel,
		permission,
		false,
		false,
		o.info.Name))

	return nil

}

func (o *Group) GetMembership(memberKey keys.UserKey) (*Membership, bool) {
	return o.memberships.GetMembership(keys.NewMembershipKey(o, memberKey))
}

func (o *Group) assertIsNew() error {
	if !o.isNew {
		return fmt.Errorf("group has already been created")
	}
	return nil
}
func (o *Group) assertNotNew() error {
	if o.isNew {
		return fmt.Errorf("group not been created yet")
	}
	return nil
}

func (o *Group) assertIsActive(key keys.UserKey) error {
	m, ok := o.GetMembership(key)
	if !ok {
		return fmt.Errorf("membership not found")
	}
	if m.Status != ApprovedMembershipStatus {
		return fmt.Errorf("member is not active")
	}
	return nil
}

func (o *Group) assertIsAdmin(userKey keys.UserKey) error {
	m, ok := o.GetMembership(userKey)
	if !ok || !m.IsAdmin() {
		return exceptions.ErrForbidden
	}
	return nil
}

func validateGroupInfo(groupInfo GroupInfo) (GroupInfo, error) {
	groupInfo = GroupInfo{
		Name:        strings.TrimSpace(groupInfo.Name),
		Description: strings.TrimSpace(groupInfo.Description),
	}

	if groupInfo.Name == "" {
		return GroupInfo{}, fmt.Errorf("group name must not be empty")
	}
	if len(groupInfo.Name) > 64 {
		return GroupInfo{}, fmt.Errorf("group name is too long")
	}
	if len(groupInfo.Name) > 1024 {
		return GroupInfo{}, fmt.Errorf("group description is too long")
	}
	return groupInfo, nil
}

func (o *Group) raise(event eventsource.Event) {
	o.changes = append(o.changes, event)
	o.on(event, true)
}

func (o *Group) on(evt eventsource.Event, isNew bool) {
	switch e := evt.(type) {
	case GroupCreated:
		o.handleGroupCreated(e)
	case GroupInfoChanged:
		o.handleGroupInfoChanged(e)
	case MembershipStatusChanged:
		o.handleMembershipStatusChange(e)
	}
	if !isNew {
		o.version++
	}
}

func (o *Group) MarkAsCommitted() {
	o.version += len(o.changes)
	o.changes = []eventsource.Event{}
}

func (o *Group) GetChanges() []eventsource.Event {
	return o.changes
}

func (o *Group) GetKey() keys.GroupKey {
	return o.key
}

func (o *Group) StreamKey() keys.StreamKey {
	return o.key.StreamKey()
}

func (o *Group) GetVersion() int {
	return o.version
}

const (
	GroupCreatedEvent            = "group_created"
	GroupInfoChangedEvent        = "group_info_changed"
	MembershipStatusChangedEvent = "group_membership_status_changed"
)

var AllEvents = []string{
	GroupCreatedEvent,
	GroupInfoChangedEvent,
	MembershipStatusChangedEvent,
}

// Group Created Event

type GroupCreatedPayload struct {
	GroupInfo GroupInfo    `json:"group_info"`
	CreatedBy keys.UserKey `json:"created_by"`
}

type GroupCreated struct {
	eventsource.EventEnvelope
	GroupCreatedPayload `json:"payload"`
}

func NewGroupCreated(createdBy keys.UserKey, groupInfo GroupInfo) GroupCreated {
	return GroupCreated{
		eventsource.NewEventEnvelope(GroupCreatedEvent, 1),
		GroupCreatedPayload{
			groupInfo,
			createdBy,
		},
	}
}

var _ eventsource.Event = GroupCreated{}

// Group Info Changed

type GroupInfoChangedPayload struct {
	NewGroupInfo GroupInfo    `json:"new_group_info"`
	OldGroupInfo GroupInfo    `json:"old_group_info"`
	ChangedBy    keys.UserKey `json:"changed_by"`
}

type GroupInfoChanged struct {
	eventsource.EventEnvelope
	GroupInfoChangedPayload `json:"payload"`
}

func NewGroupInfoChanged(changedBy keys.UserKey, oldGroupInfo, newGroupInfo GroupInfo) GroupInfoChanged {
	return GroupInfoChanged{
		eventsource.NewEventEnvelope(GroupInfoChangedEvent, 1),
		GroupInfoChangedPayload{
			newGroupInfo,
			oldGroupInfo,
			changedBy,
		},
	}
}

var _ eventsource.Event = GroupInfoChanged{}

type MembershipApproved struct {
	ApprovedBy    keys.UserKey       `json:"approved_by"`
	MembershipKey keys.MembershipKey `json:"membership_key"`
}

// Membership Status Changed

type MembershipStatusChangedPayload struct {
	OldStatus            *MembershipStatus `json:"old_status"`
	NewStatus            *MembershipStatus `json:"new_status"`
	MemberKey            keys.UserKey      `json:"member_key"`
	ChangedBy            keys.UserKey      `json:"changed_by"`
	OldPermissions       PermissionLevel   `json:"old_permissions"`
	NewPermissions       PermissionLevel   `json:"new_permissions"`
	IsNewMembership      bool              `json:"is_new_membership"`
	IsCanceledMembership bool              `json:"is_canceled_membership"`
	GroupName            string            `json:"group_name"`
}

type MembershipStatusChanged struct {
	eventsource.EventEnvelope
	MembershipStatusChangedPayload `json:"payload"`
}

func NewMembershipStatusChanged(
	changedBy keys.UserKey,
	memberKey keys.UserKey,
	oldStatus *MembershipStatus,
	newStatus *MembershipStatus,
	oldPermissions PermissionLevel,
	newPermissions PermissionLevel,
	isNewMembership bool,
	isCanceledMembership bool,
	groupName string) MembershipStatusChanged {
	return MembershipStatusChanged{
		eventsource.NewEventEnvelope(MembershipStatusChangedEvent, 1),
		MembershipStatusChangedPayload{
			oldStatus,
			newStatus,
			memberKey,
			changedBy,
			oldPermissions,
			newPermissions,
			isNewMembership,
			isCanceledMembership,
			groupName,
		},
	}
}

var _ eventsource.Event = MembershipStatusChanged{}

//

func RegisterEvents(mapper *eventsource.EventMapper) error {
	for _, eventType := range AllEvents {
		if err := mapper.RegisterMapper(eventType, MapEvent); err != nil {
			return err
		}
	}
	return nil
}

func MapEvent(eventType string, bytes []byte) (eventsource.Event, error) {
	switch eventType {
	case GroupCreatedEvent:
		var dest GroupCreated
		err := json.Unmarshal(bytes, &dest)
		return dest, err
	case GroupInfoChangedEvent:
		var dest GroupInfoChanged
		err := json.Unmarshal(bytes, &dest)
		return dest, err
	case MembershipStatusChangedEvent:
		var dest MembershipStatusChanged
		err := json.Unmarshal(bytes, &dest)
		return dest, err
	default:
		return nil, fmt.Errorf("unexpected event type '%s'", eventType)
	}
}
