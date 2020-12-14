package trading

import "github.com/commonpool/backend/model"

type OfferItem interface {
	Type() OfferItemType
	GetOfferKey() model.OfferKey
	GetKey() model.OfferItemKey
	IsCreditTransfer() bool
	IsServiceProviding() bool
	IsBorrowingResource() bool
	IsResourceTransfer() bool
	IsCompleted() bool
	IsAccepted() bool
	IsAcceptedByReceiver() bool
	IsAcceptedByGiver() bool
	GetReceiverKey() *model.Target
}
