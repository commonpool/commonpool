package domain

import "github.com/commonpool/backend/pkg/keys"

type OfferItemApproval struct {
	ApprovedBy   keys.UserKey      `json:"approved_by"`
	OfferItemKey keys.OfferItemKey `json:"offer_item_key"`
}
