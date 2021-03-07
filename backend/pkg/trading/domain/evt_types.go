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

var AllEvents = []string{
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
}

func RegisterEvents(eventMapper *eventsource.EventMapper) error {
	for _, eventType := range AllEvents {
		if err := eventMapper.RegisterMapper(eventType, MapEvent); err != nil {
			return err
		}
	}
	return nil
}

func MapEvent(eventType string, bytes []byte) (eventsource.Event, error) {

	var decoded eventsource.Event
	switch eventType {
	case OfferSubmittedEvent:
		decoded = &OfferSubmitted{}
	case OfferItemApprovedEvent:
		decoded = &OfferItemApproved{}
	case OfferDeclinedEvent:
		decoded = &OfferDeclined{}
	case OfferApprovedEvent:
		decoded = &OfferApproved{}
	case OfferCompletedEvent:
		decoded = &OfferCompleted{}
	case ServiceGivenNotifiedEvent:
		decoded = &ServiceGivenNotified{}
	case ServiceReceivedNotifiedEvent:
		decoded = &ServiceReceivedNotified{}
	case ResourceTransferGivenNotifiedEvent:
		decoded = &ResourceTransferGivenNotified{}
	case ResourceTransferReceivedNotifiedEvent:
		decoded = &ResourceTransferReceivedNotified{}
	case ResourceBorrowedNotifiedEvent:
		decoded = &ResourceBorrowedNotified{}
	case ResourceLentNotifiedEvent:
		decoded = &ResourceLentNotified{}
	case BorrowerReturnedResourceEvent:
		decoded = &BorrowerReturnedResourceNotified{}
	case LenderReceivedBackResourceEvent:
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
