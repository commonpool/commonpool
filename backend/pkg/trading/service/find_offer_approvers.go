package service

import "github.com/commonpool/backend/pkg/trading/model"

func (t TradingService) FindApproversForOffer(offerKey model.OfferKey) (*model.OfferApprovers, error) {
	return t.tradingStore.FindApproversForOffer(offerKey)
}
