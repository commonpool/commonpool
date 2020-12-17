package service

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading"
)

func (t TradingService) GetOffersForUser(key keys.UserKey) (*trading.GetOffersResult, error) {
	return t.tradingStore.GetOffersForUser(key)
}
