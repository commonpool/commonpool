package domain

type OfferEvent string

const (
	OfferSubmittedEvent                   OfferEvent = "offer_submitted"
	OfferItemApprovedEvent                OfferEvent = "offer_item_approved"
	OfferDeclinedEvent                    OfferEvent = "offer_declined"
	OfferApprovedEvent                    OfferEvent = "offer_approved"
	OfferCompletedEvent                   OfferEvent = "offer_completed"
	ServiceGivenNotifiedEvent             OfferEvent = "service_given_notified"
	ServiceReceivedNotifiedEvent          OfferEvent = "service_received_notified"
	ResourceTransferGivenNotifiedEvent    OfferEvent = "resource_transfer_given_notified"
	ResourceTransferReceivedNotifiedEvent OfferEvent = "resource_transfer_received_notified"
	ResourceBorrowedNotifiedEvent         OfferEvent = "resource_borrowed_notified"
	ResourceLentNotifiedEvent             OfferEvent = "resource_lent_notified"
	BorrowerReturnedResourceEvent         OfferEvent = "borrower_returned_resource_notified"
	LenderReceivedBackResourceEvent       OfferEvent = "lender_received_back_resource_notified"
)

type Event interface {
	GetType() OfferEvent
	GetVersion() int
}

type EventBase struct {
	key OfferEvent
}

func (e *EventBase) GetType() OfferEvent {
	return e.key
}
