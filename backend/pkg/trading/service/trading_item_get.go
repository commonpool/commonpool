package service

import (
	"context"
	"github.com/commonpool/backend/model"
	trading2 "github.com/commonpool/backend/pkg/trading"
)

func (t TradingService) GetOfferItem(ctx context.Context, offerItemKey model.OfferItemKey) (trading2.OfferItem, error) {

	offerItem, err := t.tradingStore.GetOfferItem(ctx, offerItemKey)
	if err != nil {
		return nil, err
	}

	return offerItem, nil
}
