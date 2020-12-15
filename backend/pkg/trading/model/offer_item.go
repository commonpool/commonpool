package model

import "github.com/commonpool/backend/pkg/resource/model"

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
	GetReceiverKey() *model.Target
}
