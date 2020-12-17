package transaction

import (
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource"
	"time"
)

type Entry struct {
	Key         model.TransactionEntryKey
	Type        Type
	GroupKey    keys.GroupKey
	ResourceKey *keys.ResourceKey
	Duration    *time.Duration
	Recipient   *resource.Target
	From        *resource.Target
	Timestamp   time.Time
}
