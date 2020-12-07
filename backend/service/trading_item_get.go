package service

import (
	"context"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/trading"
)

func (t TradingService) GetOfferItem(ctx context.Context, offerItemKey model.OfferItemKey) (trading.OfferItem2, error) {

	offerItem, err := t.tradingStore.GetOfferItem(ctx, offerItemKey)
	if err != nil {
		return nil, err
	}

	return offerItem, nil
}
