package trading

import (
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type OfferItemApproval struct {
	ApprovedBy   keys.UserKey      `json:"userId"`
	ApprovalSide ApprovalSide      `json:"approvalSide"`
	OfferItemKey keys.OfferItemKey `json:"offerItemKey"`
	CreatedAt    time.Time         `json:"createdAt"`
}

type OfferItemApprovalOptions struct {
	CreatedAt time.Time
}

func NewOfferItemApproval(
	offerItemKey keys.OfferItemKey,
	approvedBy keys.UserKey,
	approvalSide ApprovalSide,
	options ...OfferItemApprovalOptions) *OfferItemApproval {

	now := time.Now().UTC()
	defaultOptions := OfferItemApprovalOptions{
		CreatedAt: now,
	}
	if len(options) > 0 {
		option := options[0]
		if option.CreatedAt != time.Unix(0, 0).UTC() {
			defaultOptions.CreatedAt = option.CreatedAt
		}
	}

	return &OfferItemApproval{
		ApprovedBy:   approvedBy,
		ApprovalSide: approvalSide,
		OfferItemKey: offerItemKey,
		CreatedAt:    defaultOptions.CreatedAt,
	}
}
