package domain

import (
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type ProvideServiceItem struct {
	OfferItemBase
	ResourceKey                 keys.ResourceKey `json:"resourceId"`
	Duration                    time.Duration    `json:"duration"`
	ServiceGivenConfirmation    bool             `json:"serviceGivenConfirmation"`
	ServiceReceivedConfirmation bool             `json:"serviceReceivedConfirmation"`
	From                        *keys.Target     `json:"from"`
}

func (p *ProvideServiceItem) AsCreditTransfer() (*CreditTransferItem, bool) {
	return nil, false
}

func (p *ProvideServiceItem) AsProvideService() (*ProvideServiceItem, bool) {
	return p, true
}

func (p *ProvideServiceItem) AsBorrowResource() (*BorrowResourceItem, bool) {
	return nil, false
}

func (p *ProvideServiceItem) AsResourceTransfer() (*ResourceTransferItem, bool) {
	return nil, false
}

func (b ProvideServiceItem) GetResourceKey() keys.ResourceKey {
	return b.ResourceKey
}

func (b ProvideServiceItem) GetTo() keys.Target {
	return *b.To
}

func (p ProvideServiceItem) IsCompleted() bool {
	return p.ServiceGivenConfirmation && p.ServiceReceivedConfirmation
}

func (p ProvideServiceItem) GetType() OfferItemType {
	return ProvideService
}

var _ OfferItem = &ProvideServiceItem{}

type NewProvideServiceItemOptions struct {
	ReceiverAccepted            bool
	GiverAccepted               bool
	CreatedAt                   time.Time
	UpdatedAt                   time.Time
	ServiceGivenConfirmation    bool
	ServiceReceivedConfirmation bool
}

func NewProvideServiceItem(
	offerKey keys.OfferKey,
	offerItemKey keys.OfferItemKey,
	from *keys.Target,
	to *keys.Target,
	resourceKey keys.ResourceKey,
	duration time.Duration,
	options ...NewProvideServiceItemOptions) *ProvideServiceItem {

	now := time.Now()

	defaultOptions := &NewProvideServiceItemOptions{
		ReceiverAccepted:            false,
		GiverAccepted:               false,
		CreatedAt:                   now,
		UpdatedAt:                   now,
		ServiceGivenConfirmation:    false,
		ServiceReceivedConfirmation: false,
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
		if option.ServiceReceivedConfirmation != false {
			defaultOptions.ServiceReceivedConfirmation = option.ServiceReceivedConfirmation
		}
		if option.ServiceGivenConfirmation != false {
			defaultOptions.ServiceGivenConfirmation = option.ServiceGivenConfirmation
		}
	}

	return &ProvideServiceItem{
		OfferItemBase: OfferItemBase{
			Type:             ProvideService,
			Key:              offerItemKey,
			OfferKey:         offerKey,
			To:               to,
			ApprovedInbound:  defaultOptions.ReceiverAccepted,
			ApprovedOutbound: defaultOptions.GiverAccepted,
			CreatedAt:        defaultOptions.CreatedAt,
			UpdatedAt:        defaultOptions.UpdatedAt,
		},
		ResourceKey:                 resourceKey,
		Duration:                    duration,
		ServiceGivenConfirmation:    defaultOptions.ServiceGivenConfirmation,
		ServiceReceivedConfirmation: defaultOptions.ServiceReceivedConfirmation,
		From:                        from,
	}

}
