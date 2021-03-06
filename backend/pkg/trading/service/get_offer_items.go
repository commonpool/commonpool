package service

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
)

func (t TradingService) GetOfferItemsForOffer(offerKey keys.OfferKey) (*domain.OfferItems, error) {
	return t.tradingStore.GetOfferItemsForOffer(offerKey)
}
