package service

import (
	"context"
	"github.com/commonpool/backend/pkg/trading"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
)

func (t TradingService) GetTradingHistory(ctx context.Context, userIDs *usermodel.UserKeys) ([]trading.HistoryEntry, error) {

	tradingHistory, err := t.tradingStore.GetTradingHistory(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	return tradingHistory, nil
}
