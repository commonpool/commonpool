package transaction

import (
	"github.com/commonpool/backend/model"
	"time"
)

type Entry struct {
	Key         model.TransactionEntryKey
	Type        TransactionType
	GroupKey    model.GroupKey
	ResourceKey *model.ResourceKey
	Duration    *time.Duration
	Recipient   *model.Target
	From        *model.Target
	Timestamp   time.Time
}
