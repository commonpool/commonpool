package domain

import (
	"github.com/commonpool/backend/pkg/keys"
)

type OfferItem interface {
	Type() OfferItemType
	GetOfferKey() keys.OfferKey
	GetKey() keys.OfferItemKey
	IsCreditTransfer() bool
	IsServiceProviding() bool
	IsBorrowingResource() bool
	IsResourceTransfer() bool
	IsCompleted() bool
	IsAccepted() bool
	IsInboundApproved() bool
	IsOutboundApproved() bool
	GetReceiverKey() *Target
	AsCreditTransfer() (*CreditTransferItem, bool)
	AsProvideService() (*ProvideServiceItem, bool)
	AsBorrowResource() (*BorrowResourceItem, bool)
	AsResourceTransfer() (*ResourceTransferItem, bool)
}
