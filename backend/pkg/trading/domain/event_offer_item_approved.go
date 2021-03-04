package domain

import (
	"github.com/commonpool/backend/pkg/keys"
)

type ApprovalDirection string

const (
	Inbound  ApprovalDirection = "inbound"
	Outbound ApprovalDirection = "outbound"
)

type OfferItemApproved struct {
	ApprovedBy   keys.UserKey      `json:"approved_by"`
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
	Direction    ApprovalDirection `json:"direction"`
	Type         OfferEvent        `json:"type"`
	Version      int               `json:"version"`
}

func NewOfferItemApproved(approvedBy keys.UserKey, offerItemKey keys.OfferItemKey, direction ApprovalDirection) *OfferItemApproved {
	return &OfferItemApproved{
		ApprovedBy:   approvedBy,
		OfferItemKey: offerItemKey,
		Direction:    direction,
		Type:         OfferItemApprovedEvent,
		Version:      1,
	}
}

func (o OfferItemApproved) GetType() OfferEvent {
	return o.Type
}

func (o *OfferItemApproved) GetVersion() int {
	return o.Version
}

var _ Event = &OfferItemApproved{}
