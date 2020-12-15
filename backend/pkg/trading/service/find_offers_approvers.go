package service

import "github.com/commonpool/backend/pkg/trading/model"

func (t TradingService) FindApproversForOffers(offers *model.OfferKeys) (*model.OffersApprovers, error) {
	return t.tradingStore.FindApproversForOffers(offers)
}
