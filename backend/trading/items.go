package trading

import (
	"github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/model"
	"time"
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
		return "", errors.ErrInvalidOfferItemType
	}
}

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

type OfferItemBase struct {
	Type             OfferItemType
	Key              model.OfferItemKey
	OfferKey         model.OfferKey
	To               *model.Target
	ReceiverAccepted bool
	GiverAccepted    bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
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

func (c OfferItemBase) GetReceiverKey() *model.Target {
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
	From               *model.Target
	Amount             time.Duration
	CreditsTransferred bool
}

func (c CreditTransferItem) IsCompleted() bool {
	return c.CreditsTransferred
}

func (c CreditTransferItem) Type() OfferItemType {
	return CreditTransfer
}

var _ OfferItem = &CreditTransferItem{}

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

func (p ProvideServiceItem) Type() OfferItemType {
	return ProvideService
}

var _ OfferItem = &ProvideServiceItem{}

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

func (b BorrowResourceItem) Type() OfferItemType {
	return BorrowResource
}

var _ OfferItem = &BorrowResourceItem{}

type ResourceTransferItem struct {
	OfferItemBase
	ResourceKey  model.ResourceKey
	ItemGiven    bool
	ItemReceived bool
}

func (r ResourceTransferItem) IsCompleted() bool {
	return r.ItemGiven && r.ItemReceived
}

func (r ResourceTransferItem) Type() OfferItemType {
	return ResourceTransfer
}

var _ OfferItem = &ResourceTransferItem{}
