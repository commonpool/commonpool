package transaction

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading"
	"time"
)

type Entry struct {
	Key         keys.TransactionEntryKey
	Type        Type
	GroupKey    keys.GroupKey
	ResourceKey *keys.ResourceKey
	Duration    *time.Duration
	Recipient   *trading.Target
	From        *trading.Target
	Timestamp   time.Time
}
