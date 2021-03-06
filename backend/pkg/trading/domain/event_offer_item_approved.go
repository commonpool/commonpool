package domain

import (
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
)

type ApprovalDirection string

const (
	Inbound  ApprovalDirection = "inbound"
	Outbound ApprovalDirection = "outbound"
)

type OfferItemApprovedPayload struct {
	ApprovedBy   keys.UserKey      `json:"approved_by"`
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
	Direction    ApprovalDirection `json:"direction"`
}

type OfferItemApproved struct {
	eventsource.EventEnvelope
	OfferItemApprovedPayload `json:"payload"`
}

func NewOfferItemApproved(approvedBy keys.UserKey, offerItemKey keys.OfferItemKey, direction ApprovalDirection) *OfferItemApproved {
	return &OfferItemApproved{
		eventsource.NewEventEnvelope(OfferItemApprovedEvent, 1),
		OfferItemApprovedPayload{
			approvedBy,
			offerItemKey,
			direction,
		},
	}
}

var _ eventsource.Event = &OfferItemApproved{}
