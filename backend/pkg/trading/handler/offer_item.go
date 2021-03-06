package handler

import (
	"github.com/commonpool/backend/pkg/trading/domain"
	"time"
)

type OfferItem struct {
	ID                          string               `json:"id"`
	From                        *OfferItemTarget     `json:"from"`
	To                          *OfferItemTarget     `json:"to"`
	Type                        domain.OfferItemType `json:"type"`
	ResourceId                  *string              `json:"resourceId"`
	Duration                    *int64               `json:"duration"`
	Amount                      *int64               `json:"amount"`
	ReceiverApproved            bool                 `json:"receiverApproved"`
	GiverApproved               bool                 `json:"giverApproved"`
	ReceivingApprovers          []string             `json:"receivingApprovers"`
	GivingApprovers             []string             `json:"givingApprovers"`
	ServiceGivenConfirmation    bool                 `json:"serviceGivenConfirmation"`
	ServiceReceivedConfirmation bool                 `json:"serviceReceivedConfirmation"`
	ItemTaken                   bool                 `json:"itemTaken"`
	ItemGiven                   bool                 `json:"itemGiven"`
	ItemReturnedBack            bool                 `json:"itemReturnedBack"`
	ItemReceivedBack            bool                 `json:"itemReceivedBack"`
}

func NewResourceTransferItem(to *OfferItemTarget, resourceId string) *SendOfferPayloadItem {
	return &SendOfferPayloadItem{
		To:         *to,
		Type:       domain.ResourceTransfer,
		ResourceId: &resourceId,
	}
}

func NewCreditTransferItem(from *OfferItemTarget, to *OfferItemTarget, time time.Duration) *SendOfferPayloadItem {
	seconds := time.String()
	return &SendOfferPayloadItem{
		From:   from,
		To:     *to,
		Type:   domain.CreditTransfer,
		Amount: &seconds,
	}
}
