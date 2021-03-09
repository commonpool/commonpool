package trading

import (
	ctx "context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
)

type Service interface {
	FindTargetsForOfferItem(ctx ctx.Context, groupKey keys.GroupKey, itemType domain.OfferItemType, from *keys.Target, to *keys.Target) (*keys.Targets, error)
}
