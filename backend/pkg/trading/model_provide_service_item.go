package trading

import (
	"github.com/commonpool/backend/model"
	"time"
)

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
