package service

import (
	"github.com/commonpool/backend/pkg/trading"
)

func (t TradingService) FindApproversForOffers(offers *trading.OfferKeys) (*trading.OffersApprovers, error) {
	return t.tradingStore.FindApproversForOffers(offers)
}
