package service

import (
	"github.com/commonpool/backend/pkg/trading"
	"github.com/commonpool/backend/pkg/user/usermodel"
)

func (t TradingService) GetOffersForUser(key usermodel.UserKey) (*trading.GetOffersResult, error) {
	return t.tradingStore.GetOffersForUser(key)
}
