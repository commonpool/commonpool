package service

import "github.com/commonpool/backend/pkg/trading/model"

func (t TradingService) GetOffer(offerKey model.OfferKey) (*model.Offer, error) {
	return t.tradingStore.GetOffer(offerKey)
}
