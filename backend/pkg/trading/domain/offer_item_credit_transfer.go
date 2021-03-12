package domain

import (
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type CreditTransferItem struct {
	OfferItemBase
	From               *keys.Target  `json:"from"`
	Amount             time.Duration `json:"amount"`
	CreditsTransferred bool          `json:"creditsTransferred"`
}

func (c *CreditTransferItem) AsCreditTransfer() (*CreditTransferItem, bool) {
	return c, true
}

func (c *CreditTransferItem) AsProvideService() (*ProvideServiceItem, bool) {
	return nil, false
}

func (c *CreditTransferItem) AsBorrowResource() (*BorrowResourceItem, bool) {
	return nil, false
}

func (c *CreditTransferItem) AsResourceTransfer() (*ResourceTransferItem, bool) {
	return nil, false
}

func (c CreditTransferItem) IsCompleted() bool {
	return c.CreditsTransferred
}

func (c CreditTransferItem) GetType() OfferItemType {
	return CreditTransfer
}

var _ OfferItem = &CreditTransferItem{}

type NewCreditTransferItemOptions struct {
	ReceiverAccepted   bool
	GiverAccepted      bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
	CreditsTransferred bool
}

func NewCreditTransferItem(
	offerKey keys.OfferKey,
	offerItemKey keys.OfferItemKey,
	from *keys.Target,
	to *keys.Target,
	duration time.Duration,
	options ...NewCreditTransferItemOptions) *CreditTransferItem {

	now := time.Now()
	defaultOptions := &NewCreditTransferItemOptions{
		ReceiverAccepted:   false,
		GiverAccepted:      false,
		CreatedAt:          now,
		UpdatedAt:          now,
		CreditsTransferred: false,
	}

	if len(options) > 0 {
		option := options[0]
		if option.ReceiverAccepted != false {
			defaultOptions.ReceiverAccepted = option.ReceiverAccepted
		}
		if option.GiverAccepted != false {
			defaultOptions.GiverAccepted = option.GiverAccepted
		}
		if option.CreatedAt != time.Unix(0, 0) {
			defaultOptions.CreatedAt = option.CreatedAt
		}
		if option.UpdatedAt != time.Unix(0, 0) {
			defaultOptions.UpdatedAt = option.UpdatedAt
		}
		if option.CreditsTransferred != false {
			defaultOptions.CreditsTransferred = option.CreditsTransferred
		}
	}

	return &CreditTransferItem{
		OfferItemBase: OfferItemBase{
			Type:             CreditTransfer,
			Key:              offerItemKey,
			OfferKey:         offerKey,
			To:               to,
			ApprovedInbound:  defaultOptions.ReceiverAccepted,
			ApprovedOutbound: defaultOptions.GiverAccepted,
			CreatedAt:        defaultOptions.CreatedAt,
			UpdatedAt:        defaultOptions.UpdatedAt,
		},
		From:               from,
		Amount:             duration,
		CreditsTransferred: defaultOptions.CreditsTransferred,
	}
}
