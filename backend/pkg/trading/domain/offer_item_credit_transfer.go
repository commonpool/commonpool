package domain

import (
	"encoding/json"
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type CreditTransferItem struct {
	OfferItemKey keys.OfferItemKey `json:"key"`
	Amount       time.Duration     `json:"amount"`
	From         *OfferItemTarget  `json:"from"`
	To           *OfferItemTarget  `json:"to"`
}

func (c *CreditTransferItem) MarshalJSON() ([]byte, error) {
	a := struct {
		CreditTransferItem
		Type string `json:"type"`
	}{
		CreditTransferItem: *c,
		Type:               string(CreditTransferItemType),
	}
	return json.Marshal(a)
}

func (c CreditTransferItem) Type() OfferItemType {
	return CreditTransferItemType
}

func (c CreditTransferItem) Key() keys.OfferItemKey {
	return c.OfferItemKey
}

var _ OfferItem = &CreditTransferItem{}
