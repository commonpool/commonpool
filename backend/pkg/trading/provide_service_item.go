package trading

import (
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
	"time"
)

type ProvideServiceItem struct {
	OfferItemBase
	ResourceKey                 resourcemodel.ResourceKey
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
