package domain

import (
	"github.com/commonpool/backend/pkg/keys"
)

type OfferItemType string

const (
	CreditTransferItemType   OfferItemType = "credit_transfer"
	BorrowResourceItemType   OfferItemType = "borrow_resource"
	ResourceTransferItemType OfferItemType = "transfer_resource"
	ProvideServiceItemType   OfferItemType = "provide_service"
)

type OfferItem2 interface {
	Key() keys.OfferItemKey
	Type() OfferItemType
}
