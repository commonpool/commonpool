package service

import (
	"context"
	tradingmodel "github.com/commonpool/backend/pkg/trading/model"
)

func (t TradingService) GetOfferItem(ctx context.Context, offerItemKey tradingmodel.OfferItemKey) (tradingmodel.OfferItem, error) {

	offerItem, err := t.tradingStore.GetOfferItem(ctx, offerItemKey)
	if err != nil {
		return nil, err
	}

	return offerItem, nil
}
