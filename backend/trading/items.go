package trading

import (
	"fmt"
	"github.com/commonpool/backend/model"
	"time"
)

type OfferItemType2 string

const (
	CreditTransfer   OfferItemType2 = "transfer_credits"
	ProvideService   OfferItemType2 = "provide_service"
	BorrowResource   OfferItemType2 = "borrow_resource"
	ResourceTransfer OfferItemType2 = "transfer_resource"
)

func ParseOfferItemType(str string) (OfferItemType2, error) {
	if str == string(CreditTransfer) {
		return CreditTransfer, nil
	} else if str == string(ProvideService) {
		return ProvideService, nil
	} else if str == string(BorrowResource) {
		return BorrowResource, nil
	} else if str == string(ResourceTransfer) {
		return ResourceTransfer, nil
	} else {
		return "", fmt.Errorf("unexpected offer item type")
	}
}

type OfferItem2 interface {
	Type() OfferItemType2
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
	GetReceiverKey() *OfferItemTarget
}

type OfferItemBase struct {
	Type             OfferItemType2
	Key              model.OfferItemKey
	OfferKey         model.OfferKey
	To               *OfferItemTarget
	ReceiverAccepted bool
	GiverAccepted    bool
}

func (c OfferItemBase) GetKey() model.OfferItemKey {
	return c.Key
}

func (c OfferItemBase) GetOfferKey() model.OfferKey {
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

func (c OfferItemBase) GetReceiverKey() *OfferItemTarget {
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

type CreditTransferItem struct {
	OfferItemBase
	From               *OfferItemTarget
	Amount             time.Duration
	CreditsTransferred bool
}

func (c CreditTransferItem) IsCompleted() bool {
	return c.CreditsTransferred
}

func (c CreditTransferItem) Type() OfferItemType2 {
	return CreditTransfer
}

var _ OfferItem2 = &CreditTransferItem{}

type ProvideServiceItem struct {
	OfferItemBase
	ResourceKey                 model.ResourceKey
	Duration                    time.Duration
	ServiceGivenConfirmation    bool
	ServiceReceivedConfirmation bool
}

func (p ProvideServiceItem) IsCompleted() bool {
	return p.ServiceGivenConfirmation && p.ServiceReceivedConfirmation
}

func (p ProvideServiceItem) Type() OfferItemType2 {
	return ProvideService
}

var _ OfferItem2 = &ProvideServiceItem{}

type BorrowResourceItem struct {
	OfferItemBase
	ResourceKey      model.ResourceKey
	Duration         time.Duration
	ItemTaken        bool
	ItemGiven        bool
	ItemReturnedBack bool
	ItemReceivedBack bool
}

func (b BorrowResourceItem) IsCompleted() bool {
	return b.ItemTaken && b.ItemGiven && b.ItemReturnedBack && b.ItemReceivedBack
}

func (b BorrowResourceItem) Type() OfferItemType2 {
	return BorrowResource
}

var _ OfferItem2 = &BorrowResourceItem{}

type ResourceTransferItem struct {
	OfferItemBase
	ResourceKey  model.ResourceKey
	ItemGiven    bool
	ItemReceived bool
}

func (r ResourceTransferItem) IsCompleted() bool {
	return r.ItemGiven && r.ItemReceived
}

func (r ResourceTransferItem) Type() OfferItemType2 {
	return ResourceTransfer
}

var _ OfferItem2 = &ResourceTransferItem{}
