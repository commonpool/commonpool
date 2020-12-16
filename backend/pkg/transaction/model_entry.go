package transaction

import (
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/group"
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
	"time"
)

type Entry struct {
	Key         model.TransactionEntryKey
	Type        Type
	GroupKey    group.GroupKey
	ResourceKey *resourcemodel.ResourceKey
	Duration    *time.Duration
	Recipient   *resourcemodel.Target
	From        *resourcemodel.Target
	Timestamp   time.Time
}
