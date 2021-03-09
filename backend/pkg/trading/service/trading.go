package service

import (
	group2 "github.com/commonpool/backend/pkg/group"
	trading2 "github.com/commonpool/backend/pkg/trading"
)

type TradingService struct {
	groupService group2.Service
}

var _ trading2.Service = &TradingService{}

func NewTradingService(
	groupService group2.Service,
) *TradingService {
	return &TradingService{
		groupService: groupService,
	}
}
