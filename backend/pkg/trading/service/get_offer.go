package service

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading"
)

func (t TradingService) GetOffer(offerKey keys.OfferKey) (*trading.Offer, error) {
	return t.tradingStore.GetOffer(offerKey)
}
