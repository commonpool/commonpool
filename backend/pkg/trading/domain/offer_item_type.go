package domain

import (
	"github.com/commonpool/backend/pkg/exceptions"
)

type OfferItemType string

const (
	CreditTransfer   OfferItemType = "transfer_credits"
	ProvideService   OfferItemType = "provide_service"
	BorrowResource   OfferItemType = "borrow_resource"
	ResourceTransfer OfferItemType = "transfer_resource"
)

func ParseOfferItemType(str string) (OfferItemType, error) {
	if str == string(CreditTransfer) {
		return CreditTransfer, nil
	} else if str == string(ProvideService) {
		return ProvideService, nil
	} else if str == string(BorrowResource) {
		return BorrowResource, nil
	} else if str == string(ResourceTransfer) {
		return ResourceTransfer, nil
	} else {
		return "", exceptions.ErrInvalidOfferItemType
	}
}
