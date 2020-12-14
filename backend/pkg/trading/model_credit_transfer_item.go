package trading

import (
	"github.com/commonpool/backend/model"
	"time"
)

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
