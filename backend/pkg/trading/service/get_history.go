package service

import (
	"context"
	model2 "github.com/commonpool/backend/pkg/trading/model"
	usermodel "github.com/commonpool/backend/pkg/user/model"
)

func (t TradingService) GetTradingHistory(ctx context.Context, userIDs *usermodel.UserKeys) ([]model2.HistoryEntry, error) {

	tradingHistory, err := t.tradingStore.GetTradingHistory(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	return tradingHistory, nil
}
