package domain

import (
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/pkg/keys"
	uuid "github.com/satori/go.uuid"
)

type Offer struct {
	aggregateType                       string
	key                                 keys.OfferKey
	groupKey                            keys.GroupKey
	changes                             []Event
	version                             int
	offerItemMap                        map[keys.OfferItemKey]OfferItem
	offerItemCount                      int
	approvals                           []OfferItemApproval
	inboundApproved                     ApprovalMap
	outboundApproved                    ApprovalMap
	status                              OfferStatus
	declinedBy                          *keys.UserKey
	isNew                               bool
	serviceGivenMap                     ApprovalMap
	serviceReceivedMap                  ApprovalMap
	serviceItemMap                      map[keys.OfferItemKey]*ServiceOfferItem
	resourceTransferGivenMap            ApprovalMap
	resourceTransferReceivedMap         ApprovalMap
	resourceTransferItemMap             map[keys.OfferItemKey]*ResourceTransferItem
	creditTransferItemMap               map[keys.OfferItemKey]*CreditTransferItem
	resourceBorrowItemMap               map[keys.OfferItemKey]*ResourceBorrowItem
	resourceBorrowedMap                 ApprovalMap
	resourceLentMap                     ApprovalMap
	lenderReceivedBackBorrowedResource  ApprovalMap
	borrowerReturnedBorrowedResourceMap ApprovalMap
}

type ApprovalMap map[keys.OfferItemKey]bool

func (a *ApprovalMap) MarshalJSON() ([]byte, error) {
	val := map[string]bool{}
	for key, b := range *a {
		val[key.String()] = b
	}
	return json.Marshal(val)
}

type OfferItemMap map[keys.OfferItemKey]OfferItem

func (a *OfferItemMap) MarshalJSON() ([]byte, error) {
	val := map[string]OfferItem{}
	for key, b := range *a {
		val[key.String()] = b
	}
	return json.Marshal(val)
}

func NewOffer() *Offer {
	return &Offer{
		aggregateType:                       "offer",
		key:                                 keys.NewOfferKey(uuid.NewV4()),
		inboundApproved:                     ApprovalMap{},
		outboundApproved:                    ApprovalMap{},
		offerItemMap:                        OfferItemMap{},
		status:                              Pending,
		isNew:                               true,
		serviceGivenMap:                     ApprovalMap{},
		serviceReceivedMap:                  ApprovalMap{},
		serviceItemMap:                      map[keys.OfferItemKey]*ServiceOfferItem{},
		resourceTransferItemMap:             map[keys.OfferItemKey]*ResourceTransferItem{},
		resourceTransferGivenMap:            ApprovalMap{},
		resourceTransferReceivedMap:         ApprovalMap{},
		resourceBorrowItemMap:               map[keys.OfferItemKey]*ResourceBorrowItem{},
		resourceBorrowedMap:                 ApprovalMap{},
		resourceLentMap:                     ApprovalMap{},
		lenderReceivedBackBorrowedResource:  ApprovalMap{},
		borrowerReturnedBorrowedResourceMap: ApprovalMap{},
		creditTransferItemMap:               map[keys.OfferItemKey]*CreditTransferItem{},
		changes:                             []Event{},
	}
}

func NewFromEvents(key keys.OfferKey, events []Event) *Offer {
	offer := NewOffer()
	offer.key = key
	for _, event := range events {
		offer.on(event, false)
	}
	return offer
}

func (o *Offer) GetType() string {
	return o.aggregateType
}

// ////////////////
// //  Commands  //
// ////////////////

func (o *Offer) Submit(groupKey keys.GroupKey, offerItems OfferItems) error {

	if err := o.assertNew(); err != nil {
		return fmt.Errorf("cannot submit offer: %v", err)
	}

	if len(offerItems.Items) == 0 {
		return fmt.Errorf("cannot submit offer: must have at least one offerItem")
	}

	o.raise(NewOfferSubmitted(offerItems, groupKey))

	return nil
}

func (o *Offer) ApproveOfferItem(approver keys.UserKey, offerItemKey keys.OfferItemKey, direction ApprovalDirection, permissionMatrix PermissionMatrix) error {

	if err := o.assertNotNew(); err != nil {
		return fmt.Errorf("cannot approve offer item (%s): %v", direction, err)
	}

	if err := o.assertStatus(Pending); err != nil {
		return fmt.Errorf("cannot approve offer item (%s): %v", direction, err)
	}

	isApproved, err := o.IsOfferItemApproved(offerItemKey, direction)
	if err != nil {
		return err
	}
	if isApproved {
		return nil
	}
	offerItem := o.offerItemMap[offerItemKey]

	if !permissionMatrix(approver, offerItem, direction) {
		return fmt.Errorf("cannot approve offer item (%v): user '%s' is not allowed to do this operation", direction, approver.String())
	}

	o.raise(NewOfferItemApproved(approver, offerItemKey, direction))

	for _, item := range o.offerItemMap {
		if !o.outboundApproved[item.Key()] {
			return nil
		}
		if !o.inboundApproved[item.Key()] {
			return nil
		}
	}

	o.raise(NewOfferApproved())

	return nil
}

func (o *Offer) DeclineOffer(declinedBy keys.UserKey) error {

	if err := o.assertNotNew(); err != nil {
		return fmt.Errorf("cannot decline offer: %v", err)
	}

	if err := o.assertStatus(Pending); err != nil {
		return fmt.Errorf("cannot decline offer: %v", err)
	}

	if o.status == Declined {
		return nil
	}

	o.raise(NewOfferDeclined(declinedBy))

	return nil

}

func (o *Offer) NotifyServiceReceived(notifiedBy keys.UserKey, serviceOfferItemKey keys.OfferItemKey) error {

	if err := o.assertNotNew(); err != nil {
		return fmt.Errorf("cannot notify service received: %v", err)
	}

	if err := o.assertStatus(Approved); err != nil {
		return fmt.Errorf("cannot notify service received: %v", err)
	}

	offerItem, err := o.assertOfferItemExists(serviceOfferItemKey)
	if err != nil {
		return fmt.Errorf("cannot notify service received: %v", err)
	}

	alreadyReceived := o.serviceReceivedMap[offerItem.Key()]
	if alreadyReceived {
		return nil
	}

	serviceItem, ok := o.serviceItemMap[serviceOfferItemKey]
	if !ok {
		return fmt.Errorf("cannot notify service received: offer item is not of `ServiceOfferItem` type")
	}

	o.raise(NewServiceReceivedNotified(notifiedBy, serviceItem.Key()))

	o.CheckOfferCompleted()

	return nil

}

func (o *Offer) NotifyServiceGiven(notifiedBy keys.UserKey, serviceOfferItemKey keys.OfferItemKey) error {

	if err := o.assertNotNew(); err != nil {
		return fmt.Errorf("cannot notify service given: %v", err)
	}

	if err := o.assertStatus(Approved); err != nil {
		return fmt.Errorf("cannot notify service given: %v", err)
	}

	offerItem, err := o.assertOfferItemExists(serviceOfferItemKey)
	if err != nil {
		return fmt.Errorf("cannot notify service given: %v", err)
	}

	alreadyGiven := o.serviceGivenMap[offerItem.Key()]
	if alreadyGiven {
		return nil
	}

	serviceItem, ok := o.serviceItemMap[serviceOfferItemKey]
	if !ok {
		return fmt.Errorf("cannot notify service given: offer item is not of `ServiceOfferItem` type")
	}

	o.raise(NewServiceGivenNotified(notifiedBy, serviceItem.Key()))

	o.CheckOfferCompleted()

	return nil

}

func (o *Offer) NotifyResourceReceived(notifiedBy keys.UserKey, resourceOfferItemKey keys.OfferItemKey) error {

	if err := o.assertNotNew(); err != nil {
		return fmt.Errorf("cannot notify resource received: %v", err)
	}

	if err := o.assertStatus(Approved); err != nil {
		return fmt.Errorf("cannot notify resource received: %v", err)
	}

	offerItem, err := o.assertOfferItemExists(resourceOfferItemKey)
	if err != nil {
		return fmt.Errorf("cannot notify resource received: %v", err)
	}

	alreadyReceived := o.resourceTransferReceivedMap[offerItem.Key()]
	if alreadyReceived {
		return nil
	}

	resourceItem, ok := o.resourceTransferItemMap[resourceOfferItemKey]
	if !ok {
		return fmt.Errorf("cannot notify resource received: offer item is not of `ResourceTransferItem` type")
	}

	o.raise(NewResourceReceivedNotified(notifiedBy, resourceItem.Key()))

	o.CheckOfferCompleted()

	return nil

}

func (o *Offer) NotifyResourceGiven(notifiedBy keys.UserKey, resourceOfferItemKey keys.OfferItemKey) error {

	if err := o.assertNotNew(); err != nil {
		return fmt.Errorf("cannot notify resource given: %v", err)
	}

	if err := o.assertStatus(Approved); err != nil {
		return fmt.Errorf("cannot notify resource given: %v", err)
	}

	offerItem, err := o.assertOfferItemExists(resourceOfferItemKey)
	if err != nil {
		return fmt.Errorf("cannot notify resource given: %v", err)
	}

	alreadyGiven := o.resourceTransferGivenMap[offerItem.Key()]
	if alreadyGiven {
		return nil
	}

	resourceItem, ok := o.resourceTransferItemMap[resourceOfferItemKey]
	if !ok {
		return fmt.Errorf("cannot notify resource given: offer item is not of `ResourceTransferItem` type")
	}

	o.raise(NewResourceGivenNotified(notifiedBy, resourceItem.Key()))

	o.CheckOfferCompleted()

	return nil

}

func (o *Offer) NotifyResourceBorrowed(notifiedBy keys.UserKey, resourceOfferItemKey keys.OfferItemKey) error {

	if err := o.assertNotNew(); err != nil {
		return fmt.Errorf("cannot notify resource borrowed: %v", err)
	}

	if err := o.assertStatus(Approved); err != nil {
		return fmt.Errorf("cannot notify resource borrowed: %v", err)
	}

	offerItem, err := o.assertOfferItemExists(resourceOfferItemKey)
	if err != nil {
		return fmt.Errorf("cannot notify resource borrowed: %v", err)
	}

	alreadyBorrowed := o.resourceBorrowedMap[offerItem.Key()]
	if alreadyBorrowed {
		return nil
	}

	resourceItem, ok := o.resourceBorrowItemMap[resourceOfferItemKey]
	if !ok {
		return fmt.Errorf("cannot notify resource borrowed: offer item is not of `ResourceBorrowItem` type")
	}

	o.raise(NewResourceBorrowedNotified(notifiedBy, resourceItem.Key()))

	o.CheckOfferCompleted()

	return nil

}

func (o *Offer) NotifyResourceLent(notifiedBy keys.UserKey, resourceOfferItemKey keys.OfferItemKey) error {

	if err := o.assertNotNew(); err != nil {
		return fmt.Errorf("cannot notify resource lent: %v", err)
	}

	if err := o.assertStatus(Approved); err != nil {
		return fmt.Errorf("cannot notify resource lent: %v", err)
	}

	offerItem, err := o.assertOfferItemExists(resourceOfferItemKey)
	if err != nil {
		return fmt.Errorf("cannot notify resource lent: %v", err)
	}

	alreadyLent := o.resourceLentMap[offerItem.Key()]
	if alreadyLent {
		return nil
	}

	resourceItem, ok := o.resourceBorrowItemMap[resourceOfferItemKey]
	if !ok {
		return fmt.Errorf("cannot notify resource lent: offer item is not of `ResourceBorrowItem` type")
	}

	o.raise(NewResourceLentNotified(notifiedBy, resourceItem.Key()))

	o.CheckOfferCompleted()

	return nil

}

func (o *Offer) NotifyBorrowerReturnedResource(notifiedBy keys.UserKey, resourceOfferItemKey keys.OfferItemKey) error {

	if err := o.assertNotNew(); err != nil {
		return fmt.Errorf("cannot notify borrower returned resource: %v", err)
	}

	if err := o.assertStatus(Approved); err != nil {
		return fmt.Errorf("cannot notify borrower returned resource: %v", err)
	}

	offerItem, err := o.assertOfferItemExists(resourceOfferItemKey)
	if err != nil {
		return fmt.Errorf("cannot notify borrower returned resource: %v", err)
	}

	alreadyReturned := o.borrowerReturnedBorrowedResourceMap[offerItem.Key()]
	if alreadyReturned {
		return nil
	}

	alreadyBorrowed := o.resourceBorrowedMap[offerItem.Key()]
	if !alreadyBorrowed {
		return fmt.Errorf("cannot notify borrower returned resource: item has not been borrowed")
	}

	resourceItem, ok := o.resourceBorrowItemMap[resourceOfferItemKey]
	if !ok {
		return fmt.Errorf("cannot notify borrower returned resource: offer item is not of `ResourceBorrowItem` type")
	}

	o.raise(NewBorrowerReturnedResource(notifiedBy, resourceItem.Key()))

	o.CheckOfferCompleted()

	return nil

}

func (o *Offer) NotifyLenderReceivedBackResource(notifiedBy keys.UserKey, resourceOfferItemKey keys.OfferItemKey) error {

	if err := o.assertNotNew(); err != nil {
		return fmt.Errorf("cannot notify lender received back resource: %v", err)
	}

	if err := o.assertStatus(Approved); err != nil {
		return fmt.Errorf("cannot notify lender received back resource: %v", err)
	}

	offerItem, err := o.assertOfferItemExists(resourceOfferItemKey)
	if err != nil {
		return fmt.Errorf("cannot notify lender received back resource: %v", err)
	}

	alreadyReceivedBack := o.lenderReceivedBackBorrowedResource[offerItem.Key()]
	if alreadyReceivedBack {
		return nil
	}

	alreadyBorrowed := o.resourceLentMap[offerItem.Key()]
	if !alreadyBorrowed {
		return fmt.Errorf("cannot notify borrower returned resource: item has not been lent")
	}

	resourceItem, ok := o.resourceBorrowItemMap[resourceOfferItemKey]
	if !ok {
		return fmt.Errorf("cannot notify lender received back resource: offer item is not of `ResourceBorrowItem` type")
	}

	o.raise(NewLenderReceivedBackResource(notifiedBy, resourceItem.Key()))

	o.CheckOfferCompleted()

	return nil

}

// ///////////////
// //  Queries  //
// ///////////////

func (o *Offer) OfferItemCount() int {
	return o.offerItemCount
}

func (o *Offer) IsOfferItemApproved(offerItemKey keys.OfferItemKey, direction ApprovalDirection) (bool, error) {
	approvalMap, err := o.getApprovalMap(direction)
	if err != nil {
		return false, err
	}
	result, ok := approvalMap[offerItemKey]
	if !ok {
		return false, fmt.Errorf("offer item not part of offer")
	}
	return result, nil
}

func (o *Offer) getApprovalMap(direction ApprovalDirection) (map[keys.OfferItemKey]bool, error) {
	var approvalMap map[keys.OfferItemKey]bool
	if direction == Inbound {
		approvalMap = o.inboundApproved
	} else if direction == Outbound {
		approvalMap = o.outboundApproved
	} else {
		return nil, fmt.Errorf("invalid approval direction")
	}
	return approvalMap, nil
}

// /////////////
// //  Infra  //
// /////////////

func (o *Offer) CheckOfferCompleted() {

	// Check that all service have been given/received
	for key, _ := range o.serviceItemMap {
		if !o.serviceReceivedMap[key] {
			return
		}
		if !o.serviceGivenMap[key] {
			return
		}
	}

	for key, _ := range o.resourceTransferItemMap {
		if !o.resourceTransferGivenMap[key] {
			return
		}
		if !o.resourceTransferReceivedMap[key] {
			return
		}
	}

	for key, _ := range o.resourceBorrowItemMap {
		if !o.resourceBorrowedMap[key] {
			return
		}
		if !o.resourceLentMap[key] {
			return
		}
		if !o.borrowerReturnedBorrowedResourceMap[key] {
			return
		}
		if !o.lenderReceivedBackBorrowedResource[key] {
			return
		}
	}

	o.raise(NewOfferCompleted())

}

func (o *Offer) MarkAsCommitted() {
	o.version = o.version + len(o.changes)
	o.changes = []Event{}
}

func (o *Offer) GetChanges() []Event {
	tmp := make([]Event, len(o.changes))
	copy(tmp, o.changes)
	return tmp
}

func (o *Offer) GetKey() keys.OfferKey {
	return o.key
}

func (o *Offer) GetVersion() int {
	return o.version
}

func (o *Offer) MarshalJSON() ([]byte, error) {

	type OfferLogger struct {
		Changes          []Event             `json:"changes"`
		Version          int                 `json:"version"`
		OfferItems       []OfferItem         `json:"offer_items"`
		OfferItemMap     OfferItemMap        `json:"offer_item_map"`
		Approvals        []OfferItemApproval `json:"approvals"`
		InboundApproved  ApprovalMap         `json:"inbound_approved"`
		OutboundApproved ApprovalMap         `json:"outbound_approved"`
		Status           OfferStatus         `json:"status"`
		DeclinedBy       *keys.UserKey       `json:"declined_by"`
		IsNew            bool                `json:"is_new"`
	}

	a := OfferLogger{
		Changes:          o.changes,
		Version:          o.version,
		OfferItemMap:     o.offerItemMap,
		Approvals:        o.approvals,
		InboundApproved:  o.inboundApproved,
		OutboundApproved: o.outboundApproved,
		Status:           o.status,
		DeclinedBy:       o.declinedBy,
		IsNew:            o.isNew,
	}

	return json.Marshal(&a)
}

func (o *Offer) assertNew() error {
	if !o.isNew {
		return fmt.Errorf("offer has already been submitted")
	}
	return nil
}

func (o *Offer) assertNotNew() error {
	if o.isNew {
		return fmt.Errorf("offer has not yet been submitted")
	}
	return nil
}

func (o *Offer) assertStatus(status OfferStatus) error {
	if o.status != status {
		return fmt.Errorf("offer status must be '%s' but is '%s'", status, o.status)
	}
	return nil
}

func (o *Offer) assertOfferItemExists(offerItemKey keys.OfferItemKey) (OfferItem, error) {
	offerItem, ok := o.offerItemMap[offerItemKey]
	if !ok {
		return nil, fmt.Errorf("offer item with key %s not found", offerItemKey.String())
	}

	return offerItem, nil
}

func (o *Offer) raise(event Event) {
	o.changes = append(o.changes, event)
	o.on(event, true)
}

func (o *Offer) on(evt Event, isNew bool) {
	switch e := evt.(type) {
	case *OfferSubmitted:
		o.isNew = false
		o.groupKey = e.GroupKey
		o.offerItemCount = len(e.OfferItems.Items)
		for _, item := range e.OfferItems.Items {
			o.outboundApproved[item.Key()] = false
			o.inboundApproved[item.Key()] = false
			o.offerItemMap[item.Key()] = item

			if service, ok := item.(*ServiceOfferItem); ok {
				o.serviceItemMap[item.Key()] = service
			}
			if resourceTransfer, ok := item.(*ResourceTransferItem); ok {
				o.resourceTransferItemMap[item.Key()] = resourceTransfer
			}
			if creditTransfer, ok := item.(*CreditTransferItem); ok {
				o.creditTransferItemMap[item.Key()] = creditTransfer
			}
			if resourceBorrow, ok := item.(*ResourceBorrowItem); ok {
				o.resourceBorrowItemMap[item.Key()] = resourceBorrow
			}
		}
	case *OfferItemApproved:
		o.approvals = append(o.approvals, OfferItemApproval{
			ApprovedBy:   e.ApprovedBy,
			OfferItemKey: e.OfferItemKey,
		})
		approvalMap, _ := o.getApprovalMap(e.Direction)
		approvalMap[e.OfferItemKey] = true
	case *OfferDeclined:
		o.status = Declined
	case *OfferApproved:
		o.status = Approved
	case *OfferCompleted:
		o.status = Completed
	case *ServiceGivenNotified:
		o.serviceGivenMap[e.OfferItemKey] = true
	case *ServiceReceivedNotified:
		o.serviceReceivedMap[e.OfferItemKey] = true
	case *ResourceTransferGivenNotified:
		o.resourceTransferGivenMap[e.OfferItemKey] = true
	case *ResourceTransferReceivedNotified:
		o.resourceTransferReceivedMap[e.OfferItemKey] = true
	case *ResourceBorrowedNotified:
		o.resourceBorrowedMap[e.OfferItemKey] = true
	case *ResourceLentNotified:
		o.resourceLentMap[e.OfferItemKey] = true
	case *BorrowerReturnedResourceNotified:
		o.borrowerReturnedBorrowedResourceMap[e.OfferItemKey] = true
	case *LenderReceivedBackResourceNotified:
		o.lenderReceivedBackBorrowedResource[e.OfferItemKey] = true
	}

	if !isNew {
		o.version++
	}
}
