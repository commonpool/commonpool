package trading

import "github.com/commonpool/backend/pkg/keys"

type Approver struct {
	UserKey      keys.UserKey      `json:"userId"`
	OfferItemKey keys.OfferItemKey `json:"offerItemId"`
	ApprovalSide ApprovalSide      `json:"approvalSide"`
}

func NewApprover(userKey keys.UserKey, offerItemKey keys.OfferItemKey, approvalSide ApprovalSide) *Approver {
	return &Approver{
		UserKey:      userKey,
		OfferItemKey: offerItemKey,
		ApprovalSide: approvalSide,
	}
}
