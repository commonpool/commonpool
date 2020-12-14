package service

import (
	"context"
	"github.com/commonpool/backend/model"
	trading2 "github.com/commonpool/backend/pkg/trading"
)

func (t TradingService) GetTradingHistory(ctx context.Context, userIDs *model.UserKeys) ([]trading2.HistoryEntry, error) {

	tradingHistory, err := t.tradingStore.GetTradingHistory(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	return tradingHistory, nil
}
