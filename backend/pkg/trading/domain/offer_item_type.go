package domain

import (
	"encoding/json"
	"fmt"
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

func (o OfferItemType) GormDataType() string {
	return "varchar(64)"
}

func (o *OfferItemType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	switch str {
	case string(CreditTransfer):
		*o = CreditTransfer
	case string(ProvideService):
		*o = ProvideService
	case string(BorrowResource):
		*o = BorrowResource
	case string(ResourceTransfer):
		*o = ResourceTransfer
	default:
		return fmt.Errorf("invalid offer item type : %s", str)
	}
	return nil
}
