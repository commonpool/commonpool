package trading

import (
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/keys"
)

type OfferApprovalState struct {
	Offer      *Offer              `json:"offer"`
	OfferItems *OfferItems         `json:"offerItems"`
	Approvals  *OfferItemApprovals `json:"approvals"`
}

func (s OfferApprovalState) GetOfferItem(key keys.OfferItemKey) (OfferItem, error) {
	for _, item := range s.OfferItems.Items {
		if item.GetKey() == key {
			return item, nil
		}
	}
	return nil, exceptions.ErrOfferItemNotFound
}

func (s OfferApprovalState) IsOutboundApproved(key keys.OfferItemKey) (bool, error) {
	return s.IsApproved(key, Outbound)
}

func (s OfferApprovalState) IsInboundApproved(key keys.OfferItemKey) (bool, error) {
	return s.IsApproved(key, Inbound)
}

func (s OfferApprovalState) IsApproved(key keys.OfferItemKey, direction ApprovalSide) (bool, error) {
	var offerItemFound = false
	for _, approval := range s.Approvals.Items {
		if approval.OfferItemKey == key {
			offerItemFound = true
			if approval.ApprovalSide == direction {
				return true, nil
			}
		}
	}
	if !offerItemFound {
		return false, exceptions.ErrOfferItemNotFound
	}
	return false, nil
}

func (s OfferApprovalState) GetApprovalsForOfferItem(key keys.OfferItemKey) *OfferItemApprovals {
	var approvals []*OfferItemApproval
	for _, approval := range s.Approvals.Items {
		if approval.OfferItemKey == key {
			approvals = append(approvals, approval)
		}
	}
	return NewOfferItemApprovals(approvals...)
}

func NewOfferApprovalState(offer *Offer, offerItems *OfferItems, approvals *OfferItemApprovals) *OfferApprovalState {
	return &OfferApprovalState{
		Offer:      offer,
		OfferItems: offerItems,
		Approvals:  approvals,
	}
}
