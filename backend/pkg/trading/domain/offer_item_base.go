package domain

import (
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type OfferItemBase struct {
	Type             OfferItemType     `json:"type"`
	Key              keys.OfferItemKey `json:"id"`
	OfferKey         keys.OfferKey     `json:"offerId"`
	To               *keys.Target      `json:"to"`
	ApprovedInbound  bool              `json:"approvedInbound"`
	ApprovedOutbound bool              `json:"approvedOutbound"`
	CreatedAt        time.Time         `json:"createdAt"`
	UpdatedAt        time.Time         `json:"updatedAt"`
}

func (c OfferItemBase) GetKey() keys.OfferItemKey {
	return c.Key
}

func (c OfferItemBase) GetOfferKey() keys.OfferKey {
	return c.OfferKey
}

func (c OfferItemBase) IsCreditTransfer() bool {
	return c.Type == CreditTransfer
}
func (c OfferItemBase) IsServiceProviding() bool {
	return c.Type == ProvideService
}

func (c OfferItemBase) IsBorrowingResource() bool {
	return c.Type == BorrowResource
}

func (c OfferItemBase) IsResourceTransfer() bool {
	return c.Type == ResourceTransfer
}

func (c OfferItemBase) GetReceiverKey() *keys.Target {
	return c.To
}

func (c OfferItemBase) IsAccepted() bool {
	return c.IsOutboundApproved() && c.IsInboundApproved()
}

func (c OfferItemBase) IsInboundApproved() bool {
	return c.ApprovedInbound
}

func (c OfferItemBase) IsOutboundApproved() bool {
	return c.ApprovedOutbound
}
