package transaction

import (
	"github.com/commonpool/backend/model"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
	"time"
)

type Entry struct {
	Key         model.TransactionEntryKey
	Type        Type
	GroupKey    groupmodel.GroupKey
	ResourceKey *resourcemodel.ResourceKey
	Duration    *time.Duration
	Recipient   *model.Target
	From        *model.Target
	Timestamp   time.Time
}
