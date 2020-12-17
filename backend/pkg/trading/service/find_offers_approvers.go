package service

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading"
)

func (t TradingService) FindApproversForOffers(offers *keys.OfferKeys) (*trading.OffersApprovers, error) {
	return t.tradingStore.FindApproversForOffers(offers)
}
