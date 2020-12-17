package service

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading"
)

func (t TradingService) GetOfferItemsForOffer(offerKey keys.OfferKey) (*trading.OfferItems, error) {
	return t.tradingStore.GetOfferItemsForOffer(offerKey)
}
