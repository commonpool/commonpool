package service

import (
	"github.com/commonpool/backend/pkg/trading"
	"github.com/commonpool/backend/pkg/user/model"
)

func (t TradingService) GetOffersForUser(key model.UserKey) (*trading.GetOffersResult, error) {
	return t.tradingStore.GetOffersForUser(key)
}
