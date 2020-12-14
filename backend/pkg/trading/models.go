package trading

import (
	"github.com/commonpool/backend/model"
	"time"
)

type HistoryEntry struct {
	Timestamp         time.Time
	FromUserID        model.UserKey
	ToUserID          model.UserKey
	ResourceID        *model.ResourceKey
	TimeAmountSeconds *int64
}
