package trading

import (
	"github.com/commonpool/backend/pkg/resource"
)

type OfferItem interface {
	Type() OfferItemType
	GetOfferKey() OfferKey
	GetKey() OfferItemKey
	IsCreditTransfer() bool
	IsServiceProviding() bool
	IsBorrowingResource() bool
	IsResourceTransfer() bool
	IsCompleted() bool
	IsAccepted() bool
	IsAcceptedByReceiver() bool
	IsAcceptedByGiver() bool
	GetReceiverKey() *resource.Target
}
