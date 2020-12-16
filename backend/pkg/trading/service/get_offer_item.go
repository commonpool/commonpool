package service

import (
	"context"
	"github.com/commonpool/backend/pkg/trading"
)

func (t TradingService) GetOfferItem(ctx context.Context, offerItemKey trading.OfferItemKey) (trading.OfferItem, error) {

	offerItem, err := t.tradingStore.GetOfferItem(ctx, offerItemKey)
	if err != nil {
		return nil, err
	}

	return offerItem, nil
}
