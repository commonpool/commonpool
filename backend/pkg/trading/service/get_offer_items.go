package service

import "github.com/commonpool/backend/pkg/trading/model"

func (t TradingService) GetOfferItemsForOffer(offerKey model.OfferKey) (*model.OfferItems, error) {
	return t.tradingStore.GetOfferItemsForOffer(offerKey)
}
