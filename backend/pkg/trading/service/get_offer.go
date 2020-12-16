package service

import (
	"github.com/commonpool/backend/pkg/trading"
)

func (t TradingService) GetOffer(offerKey trading.OfferKey) (*trading.Offer, error) {
	return t.tradingStore.GetOffer(offerKey)
}
