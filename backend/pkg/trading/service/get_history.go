package service

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading"
)

func (t TradingService) GetTradingHistory(ctx context.Context, userIDs *keys.UserKeys) ([]trading.HistoryEntry, error) {

	tradingHistory, err := t.tradingStore.GetTradingHistory(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	return tradingHistory, nil
}
