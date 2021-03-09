package domain

import (
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/keys"
)

type Offer struct {
	aggregateType                       string
	key                                 keys.OfferKey
	groupKey                            keys.GroupKey
	changes                             []eventsource.Event
	version                             int
	offerItemMap                        map[keys.OfferItemKey]OfferItem
	offerItems                          *OfferItems
	offerItemCount                      int
	approvals                           []OfferItemApproval
	inboundApproved                     ApprovalMap
	outboundApproved                    ApprovalMap
	status                              OfferStatus
	declinedBy                          *keys.UserKey
	isNew                               bool
	serviceGivenMap                     ApprovalMap
	serviceReceivedMap                  ApprovalMap
	serviceItemMap                      map[keys.OfferItemKey]*ProvideServiceItem
	resourceTransferGivenMap            ApprovalMap
	resourceTransferReceivedMap         ApprovalMap
	resourceTransferItemMap             map[keys.OfferItemKey]*ResourceTransferItem
	creditTransferItemMap               map[keys.OfferItemKey]*CreditTransferItem
	resourceBorrowItemMap               map[keys.OfferItemKey]*BorrowResourceItem
	resourceBorrowedMap                 ApprovalMap
	resourceLentMap                     ApprovalMap
	lenderReceivedBackBorrowedResource  ApprovalMap
	borrowerReturnedBorrowedResourceMap ApprovalMap
	SubmittedBy                         keys.UserKey
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

func NewOffer(key keys.OfferKey) *Offer {
	return &Offer{
		aggregateType:                       "offer",
		key:                                 key,
		inboundApproved:                     ApprovalMap{},
		outboundApproved:                    ApprovalMap{},
		offerItemMap:                        OfferItemMap{},
		status:                              Pending,
		isNew:                               true,
		serviceGivenMap:                     ApprovalMap{},
		serviceReceivedMap:                  ApprovalMap{},
		serviceItemMap:                      map[keys.OfferItemKey]*ProvideServiceItem{},
		resourceTransferItemMap:             map[keys.OfferItemKey]*ResourceTransferItem{},
		resourceTransferGivenMap:            ApprovalMap{},
		resourceTransferReceivedMap:         ApprovalMap{},
		resourceBorrowItemMap:               map[keys.OfferItemKey]*BorrowResourceItem{},
		resourceBorrowedMap:                 ApprovalMap{},
		resourceLentMap:                     ApprovalMap{},
		lenderReceivedBackBorrowedResource:  ApprovalMap{},
		borrowerReturnedBorrowedResourceMap: ApprovalMap{},
		creditTransferItemMap:               map[keys.OfferItemKey]*CreditTransferItem{},
		changes:                             []eventsource.Event{},
		offerItems:                          NewOfferItems([]OfferItem{}),
	}
}

func NewFromEvents(key keys.OfferKey, events []eventsource.Event) *Offer {
	offer := NewOffer(key)
	for _, event := range events {
		offer.on(event, false)
	}
	return offer
}

func (o *Offer) GetEventType() string {
	return o.aggregateType
}

// ////////////////
// //  Commands  //
// ////////////////

func (o *Offer) Submit(submittedBy keys.UserKey, groupKey keys.GroupKey, offerItems SubmitOfferItems) error {

	if err := o.assertNew(); err != nil {
		return fmt.Errorf("cannot submit offer: %v", err)
	}

	if len(offerItems) == 0 {
		return fmt.Errorf("cannot submit offer: must have at least one offerItem")
	}

	var mappedItems = NewEmptyOfferItems()
	for _, offerItem := range offerItems {

		var item OfferItem
		switch offerItem.OfferItemType {
		case ResourceTransfer:
			if offerItem.ResourceKey == nil {
				return exceptions.ErrBadRequestf("OfferItem.ResourceKey is required for OfferItem type ResourceTransfer")
			}
			if offerItem.To == nil {
				return exceptions.ErrBadRequest("OfferItem.To is required for OfferItem type ResourceTransfer")
			}
			if offerItem.From != nil {
				return exceptions.ErrBadRequest("OfferItem.From is not valid for OfferItem type ResourceTransfer")
			}
			if offerItem.Duration != nil {
				return exceptions.ErrBadRequest("OfferItem.Duration is not valid for OfferItem type ResourceTransfer")
			}
			if offerItem.Amount != nil {
				return exceptions.ErrBadRequest("OfferItem.Amount is not valid for OfferItem type ResourceTransfer")
			}
			item = NewResourceTransferItem(o.key, offerItem.OfferItemKey, offerItem.To, *offerItem.ResourceKey)
		case BorrowResource:
			if offerItem.ResourceKey == nil {
				return exceptions.ErrBadRequestf("OfferItem.ResourceKey is required for OfferItem type BorrowResource")
			}
			if offerItem.To == nil {
				return exceptions.ErrBadRequest("OfferItem.To is required for OfferItem type BorrowResource")
			}
			if offerItem.From != nil {
				return exceptions.ErrBadRequest("OfferItem.From is not valid for OfferItem type BorrowResource")
			}
			if offerItem.Duration == nil {
				return exceptions.ErrBadRequest("OfferItem.Duration is required for OfferItem type BorrowResource")
			}
			if *offerItem.Duration < 0 {
				return exceptions.ErrBadRequest("OfferItem.Duration must be positive")
			}
			if offerItem.Amount != nil {
				return exceptions.ErrBadRequest("OfferItem.Amount is not valid for OfferItem type BorrowResource")
			}
			item = NewBorrowResourceItem(o.key, offerItem.OfferItemKey, *offerItem.ResourceKey, offerItem.To, *offerItem.Duration)
		case ProvideService:
			if offerItem.ResourceKey == nil {
				return exceptions.ErrBadRequestf("OfferItem.ResourceKey is required for OfferItem type ProvideService")
			}
			if offerItem.To == nil {
				return exceptions.ErrBadRequest("OfferItem.To is required for OfferItem type ProvideService")
			}
			if offerItem.From == nil {
				return exceptions.ErrBadRequest("OfferItem.From is required for OfferItem type ProvideService")
			}
			if offerItem.Duration == nil {
				return exceptions.ErrBadRequest("OfferItem.Duration is required for OfferItem type ProvideService")
			}
			if *offerItem.Duration < 0 {
				return exceptions.ErrBadRequest("OfferItem.Duration must be positive")
			}
			if offerItem.Amount != nil {
				return exceptions.ErrBadRequest("OfferItem.Amount is not valid for OfferItem type ProvideService")
			}
			item = NewProvideServiceItem(o.key, offerItem.OfferItemKey, offerItem.From, offerItem.To, *offerItem.ResourceKey, *offerItem.Duration)
		case CreditTransfer:
			if offerItem.ResourceKey != nil {
				return exceptions.ErrBadRequestf("OfferItem.ResourceKey is not valid for OfferItem type CreditTransfer")
			}
			if offerItem.To == nil {
				return exceptions.ErrBadRequest("OfferItem.To is required for OfferItem type CreditTransfer")
			}
			if offerItem.From == nil {
				return exceptions.ErrBadRequest("OfferItem.From is required for OfferItem type CreditTransfer")
			}
			if offerItem.Duration != nil {
				return exceptions.ErrBadRequest("OfferItem.Duration is not valid for OfferItem type CreditTransfer")
			}
			if offerItem.Amount == nil {
				return exceptions.ErrBadRequest("OfferItem.Amount is required for OfferItem type CreditTransfer")
			}
			if *offerItem.Amount < 0 {
				return exceptions.ErrBadRequest("OfferItem.Amount must be positive")
			}
			item = NewCreditTransferItem(o.key, offerItem.OfferItemKey, offerItem.From, offerItem.To, *offerItem.Amount)
		default:
			return exceptions.ErrBadRequestf("invalid offer item type: %s", offerItem.OfferItemType)
		}

		mappedItems.Append(item)

	}

	o.raise(NewOfferSubmitted(submittedBy, mappedItems, groupKey))

	return nil
}

func (o *Offer) ApproveAll(approver keys.UserKey, permissionGetter OfferPermissionGetter) error {
	for _, direction := range []ApprovalDirection{Inbound, Outbound} {
		for offerItemKey, item := range o.offerItemMap {
			if o.status == Approved {
				return nil
			}
			if permissionGetter.Can(approver, item, direction) {
				if err := o.ApproveOfferItem(approver, offerItemKey, direction, permissionGetter); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (o *Offer) ApproveOfferItem(approver keys.UserKey, offerItemKey keys.OfferItemKey, direction ApprovalDirection, permissionMatrix OfferPermissionGetter) error {

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

	if !permissionMatrix.Can(approver, offerItem, direction) {
		return fmt.Errorf("cannot approve offer item (%v): user '%s' is not allowed to do this operation", direction, approver.String())
	}

	o.raise(NewOfferItemApproved(approver, offerItemKey, direction))

	for _, item := range o.offerItemMap {
		if !o.outboundApproved[item.GetKey()] {
			return nil
		}
		if !o.inboundApproved[item.GetKey()] {
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

	alreadyReceived := o.serviceReceivedMap[offerItem.GetKey()]
	if alreadyReceived {
		return nil
	}

	serviceItem, ok := o.serviceItemMap[serviceOfferItemKey]
	if !ok {
		return fmt.Errorf("cannot notify service received: offer item is not of `ServiceOfferItem` type")
	}

	o.raise(NewServiceReceivedNotified(notifiedBy, serviceItem.GetKey()))

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

	alreadyGiven := o.serviceGivenMap[offerItem.GetKey()]
	if alreadyGiven {
		return nil
	}

	serviceItem, ok := o.serviceItemMap[serviceOfferItemKey]
	if !ok {
		return fmt.Errorf("cannot notify service given: offer item is not of `ServiceOfferItem` type")
	}

	o.raise(NewServiceGivenNotified(notifiedBy, serviceItem.GetKey()))

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

	alreadyReceived := o.resourceTransferReceivedMap[offerItem.GetKey()]
	if alreadyReceived {
		return nil
	}

	resourceItem, ok := o.resourceTransferItemMap[resourceOfferItemKey]
	if !ok {
		return fmt.Errorf("cannot notify resource received: offer item is not of `ResourceTransferItem` type")
	}

	o.raise(NewResourceReceivedNotified(notifiedBy, resourceItem.GetKey()))

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

	alreadyGiven := o.resourceTransferGivenMap[offerItem.GetKey()]
	if alreadyGiven {
		return nil
	}

	resourceItem, ok := o.resourceTransferItemMap[resourceOfferItemKey]
	if !ok {
		return fmt.Errorf("cannot notify resource given: offer item is not of `ResourceTransferItem` type")
	}

	o.raise(NewResourceGivenNotified(notifiedBy, resourceItem.GetKey()))

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

	alreadyBorrowed := o.resourceBorrowedMap[offerItem.GetKey()]
	if alreadyBorrowed {
		return nil
	}

	resourceItem, ok := o.resourceBorrowItemMap[resourceOfferItemKey]
	if !ok {
		return fmt.Errorf("cannot notify resource borrowed: offer item is not of `ResourceBorrowItem` type")
	}

	o.raise(NewResourceBorrowedNotified(notifiedBy, resourceItem.GetKey()))

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

	alreadyLent := o.resourceLentMap[offerItem.GetKey()]
	if alreadyLent {
		return nil
	}

	resourceItem, ok := o.resourceBorrowItemMap[resourceOfferItemKey]
	if !ok {
		return fmt.Errorf("cannot notify resource lent: offer item is not of `ResourceBorrowItem` type")
	}

	o.raise(NewResourceLentNotified(notifiedBy, resourceItem.GetKey()))

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

	alreadyReturned := o.borrowerReturnedBorrowedResourceMap[offerItem.GetKey()]
	if alreadyReturned {
		return nil
	}

	alreadyBorrowed := o.resourceBorrowedMap[offerItem.GetKey()]
	if !alreadyBorrowed {
		return fmt.Errorf("cannot notify borrower returned resource: item has not been borrowed")
	}

	resourceItem, ok := o.resourceBorrowItemMap[resourceOfferItemKey]
	if !ok {
		return fmt.Errorf("cannot notify borrower returned resource: offer item is not of `ResourceBorrowItem` type")
	}

	o.raise(NewBorrowerReturnedResource(notifiedBy, resourceItem.GetKey()))

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

	alreadyReceivedBack := o.lenderReceivedBackBorrowedResource[offerItem.GetKey()]
	if alreadyReceivedBack {
		return nil
	}

	alreadyBorrowed := o.resourceLentMap[offerItem.GetKey()]
	if !alreadyBorrowed {
		return fmt.Errorf("cannot notify borrower returned resource: item has not been lent")
	}

	resourceItem, ok := o.resourceBorrowItemMap[resourceOfferItemKey]
	if !ok {
		return fmt.Errorf("cannot notify lender received back resource: offer item is not of `ResourceBorrowItem` type")
	}

	o.raise(NewLenderReceivedBackResource(notifiedBy, resourceItem.GetKey()))

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

func (o *Offer) StreamKey() keys.StreamKey {
	return o.key.StreamKey()
}

func (o *Offer) MarkAsCommitted() {
	o.version = o.version + len(o.changes)
	o.changes = []eventsource.Event{}
}

func (o *Offer) GetChanges() []eventsource.Event {
	tmp := make([]eventsource.Event, len(o.changes))
	copy(tmp, o.changes)
	return tmp
}

func (o *Offer) GetOfferItems() *OfferItems {
	return o.offerItems
}

func (o *Offer) GetResourceKeys() *keys.ResourceKeys {
	var result []keys.ResourceKey
	for _, offerItem := range o.offerItems.Items {
		if offerItem.IsBorrowingResource() {
			borrowResource, _ := offerItem.AsBorrowResource()
			result = append(result, borrowResource.ResourceKey)
		} else if offerItem.IsServiceProviding() {
			serviceProviding, _ := offerItem.AsProvideService()
			result = append(result, serviceProviding.ResourceKey)
		} else if offerItem.IsResourceTransfer() {
			resourceTransfer, _ := offerItem.AsResourceTransfer()
			result = append(result, resourceTransfer.ResourceKey)
		}
	}
	return keys.NewResourceKeys(result)
}

func (o *Offer) GetKey() keys.OfferKey {
	return o.key
}

func (o *Offer) GetGroupKey() keys.GroupKey {
	return o.groupKey
}

func (o *Offer) GetSequenceNo() int {
	return o.version
}

func (o *Offer) GetVersion() int {
	return o.version
}

func (o *Offer) MarshalJSON() ([]byte, error) {

	type OfferLogger struct {
		Changes          []eventsource.Event `json:"changes"`
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

func (o *Offer) raise(event eventsource.Event) {
	o.changes = append(o.changes, event)
	o.on(event, true)
}

func (o *Offer) on(evt eventsource.Event, isNew bool) {
	switch e := evt.(type) {
	case *OfferSubmitted:
		o.isNew = false
		o.groupKey = e.GroupKey
		o.offerItemCount = len(e.OfferItems.Items)
		o.SubmittedBy = e.SubmittedBy
		for _, item := range e.OfferItems.Items {

			o.outboundApproved[item.GetKey()] = false
			o.inboundApproved[item.GetKey()] = false
			o.offerItemMap[item.GetKey()] = item
			o.offerItems.Append(item)

			if service, ok := item.AsProvideService(); ok {
				o.serviceItemMap[service.GetKey()] = service
			}

			if resourceTransfer, ok := item.AsResourceTransfer(); ok {
				o.resourceTransferItemMap[resourceTransfer.GetKey()] = resourceTransfer
			}

			if creditTransfer, ok := item.AsCreditTransfer(); ok {
				o.creditTransferItemMap[creditTransfer.GetKey()] = creditTransfer
			}

			if resourceBorrow, ok := item.AsBorrowResource(); ok {
				o.resourceBorrowItemMap[resourceBorrow.GetKey()] = resourceBorrow
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
