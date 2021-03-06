package transaction

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
	"time"
)

type Entry struct {
	Key         keys.TransactionEntryKey
	Type        Type
	GroupKey    keys.GroupKey
	ResourceKey *keys.ResourceKey
	Duration    *time.Duration
	Recipient   *domain.Target
	From        *domain.Target
	Timestamp   time.Time
}
