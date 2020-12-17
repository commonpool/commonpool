package trading

import (
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type HistoryEntry struct {
	Timestamp         time.Time
	FromUserID        keys.UserKey
	ToUserID          keys.UserKey
	ResourceID        *keys.ResourceKey
	TimeAmountSeconds *int64
}
