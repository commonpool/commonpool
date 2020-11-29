package service

import (
	"context"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/trading"
)

func (t TradingService) GetTradingHistory(ctx context.Context, userIDs *model.UserKeys) ([]trading.TradingHistoryEntry, error) {

	tradingHistory, err := t.tradingStore.GetTradingHistory(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	return tradingHistory, nil
}
