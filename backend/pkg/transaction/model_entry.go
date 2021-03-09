package transaction

import (
	"github.com/commonpool/backend/pkg/keys"
	"time"
)

type Entry struct {
	Key         keys.TransactionEntryKey
	Type        Type
	GroupKey    keys.GroupKey
	ResourceKey *keys.ResourceKey
	Duration    *time.Duration
	Recipient   *keys.Target
	From        *keys.Target
	Timestamp   time.Time
}
