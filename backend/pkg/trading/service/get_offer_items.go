package service

import (
	"github.com/commonpool/backend/pkg/trading"
)

func (t TradingService) GetOfferItemsForOffer(offerKey trading.OfferKey) (*trading.OfferItems, error) {
	return t.tradingStore.GetOfferItemsForOffer(offerKey)
}
