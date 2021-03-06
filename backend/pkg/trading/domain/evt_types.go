package domain

import (
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/pkg/eventsource"
)

const (
	OfferSubmittedEvent                   = "offer_submitted"
	OfferItemApprovedEvent                = "offer_item_approved"
	OfferDeclinedEvent                    = "offer_declined"
	OfferApprovedEvent                    = "offer_approved"
	OfferCompletedEvent                   = "offer_completed"
	ServiceGivenNotifiedEvent             = "service_given_notified"
	ServiceReceivedNotifiedEvent          = "service_received_notified"
	ResourceTransferGivenNotifiedEvent    = "resource_transfer_given_notified"
	ResourceTransferReceivedNotifiedEvent = "resource_transfer_received_notified"
	ResourceBorrowedNotifiedEvent         = "resource_borrowed_notified"
	ResourceLentNotifiedEvent             = "resource_lent_notified"
	BorrowerReturnedResourceEvent         = "borrower_returned_resource_notified"
	LenderReceivedBackResourceEvent       = "lender_received_back_resource_notified"
)

func RegisterEvents(eventMapper *eventsource.EventMapper) error {
	for _, eventType := range []string{
		OfferSubmittedEvent,
		OfferItemApprovedEvent,
		OfferDeclinedEvent,
		OfferApprovedEvent,
		OfferCompletedEvent,
		ServiceGivenNotifiedEvent,
		ServiceReceivedNotifiedEvent,
		ResourceTransferGivenNotifiedEvent,
		ResourceTransferReceivedNotifiedEvent,
		ResourceBorrowedNotifiedEvent,
		ResourceLentNotifiedEvent,
		BorrowerReturnedResourceEvent,
		LenderReceivedBackResourceEvent,
	} {
		if err := eventMapper.RegisterMapper(eventType, MapEvent); err != nil {
			return err
		}
	}
	return nil
}

func MapEvent(eventType string, bytes []byte) (eventsource.Event, error) {

	var decoded eventsource.Event
	switch eventType {
	case string(OfferSubmittedEvent):
		decoded = &OfferSubmitted{}
	case string(OfferItemApprovedEvent):
		decoded = &OfferItemApproved{}
	case string(OfferDeclinedEvent):
		decoded = &OfferDeclined{}
	case string(OfferApprovedEvent):
		decoded = &OfferApproved{}
	case string(OfferCompletedEvent):
		decoded = &OfferCompleted{}
	case string(ServiceGivenNotifiedEvent):
		decoded = &ServiceGivenNotified{}
	case string(ServiceReceivedNotifiedEvent):
		decoded = &ServiceReceivedNotified{}
	case string(ResourceTransferGivenNotifiedEvent):
		decoded = &ResourceTransferGivenNotified{}
	case string(ResourceTransferReceivedNotifiedEvent):
		decoded = &ResourceTransferReceivedNotified{}
	case string(ResourceBorrowedNotifiedEvent):
		decoded = &ResourceBorrowedNotified{}
	case string(ResourceLentNotifiedEvent):
		decoded = &ResourceLentNotified{}
	case string(BorrowerReturnedResourceEvent):
		decoded = &BorrowerReturnedResourceNotified{}
	case string(LenderReceivedBackResourceEvent):
		decoded = &LenderReceivedBackResourceNotified{}
	default:
		return nil, fmt.Errorf("unexpected event type: %s", eventType)
	}

	err := json.Unmarshal(bytes, decoded)
	if err != nil {
		return nil, err
	}

	return decoded, nil

}
