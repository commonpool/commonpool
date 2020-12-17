package service

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading"
)

func (t TradingService) FindApproversForOffer(offerKey keys.OfferKey) (*trading.OfferApprovers, error) {
	return t.tradingStore.FindApproversForOffer(offerKey)
}
