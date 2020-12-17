package trading

import (
	"github.com/commonpool/backend/pkg/resource"
	"time"
)

type OfferItemBase struct {
	Type             OfferItemType
	Key              OfferItemKey
	OfferKey         OfferKey
	To               *resource.Target
	ReceiverAccepted bool
	GiverAccepted    bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (c OfferItemBase) GetKey() OfferItemKey {
	return c.Key
}

func (c OfferItemBase) GetOfferKey() OfferKey {
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

func (c OfferItemBase) GetReceiverKey() *resource.Target {
	return c.To
}

func (c OfferItemBase) IsAccepted() bool {
	return c.IsAcceptedByGiver() && c.IsAcceptedByReceiver()
}

func (c OfferItemBase) IsAcceptedByReceiver() bool {
	return c.ReceiverAccepted
}

func (c OfferItemBase) IsAcceptedByGiver() bool {
	return c.GiverAccepted
}
