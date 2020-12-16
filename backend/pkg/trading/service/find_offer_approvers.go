package service

import (
	"github.com/commonpool/backend/pkg/trading"
)

func (t TradingService) FindApproversForOffer(offerKey trading.OfferKey) (*trading.OfferApprovers, error) {
	return t.tradingStore.FindApproversForOffer(offerKey)
}
