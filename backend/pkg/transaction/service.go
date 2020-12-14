package transaction

import (
	"github.com/commonpool/backend/model"
	"time"
)

type Service interface {
	UserSharedResourceWithGroup(groupKey model.GroupKey, resourceKey model.ResourceKey) (*Entry, error)
	UserRemovedResourceFromGroup(groupKey model.GroupKey, resourceKey model.ResourceKey) (*Entry, error)
	ServiceWasProvided(groupKey model.GroupKey, resourceKey model.ResourceKey, duration time.Duration) (*Entry, error)
	ResourceWasBorrowed(groupKey model.GroupKey, resourceKey model.ResourceKey, recipient *model.Target, expectedDuration time.Duration) (*Entry, error)
	ResourceWasReturned(groupKey model.GroupKey, resourceKey model.ResourceKey, recipient *model.Target, actualDuration time.Duration) (*Entry, error)
	ResourceWasTaken(groupKey model.GroupKey, resourceKey model.ResourceKey, recipient *model.Target) (*Entry, error)
	TimeCreditsExchanged(groupKey model.GroupKey, from *model.Target, recipient *model.Target, amount time.Duration) (*Entry, error)
	GetEntriesForGroupAndUsers(groupKey model.GroupKey, userKeys *model.UserKeys) (*Entries, error)
	GetEntry(entryKey model.TransactionEntryKey) (*Entry, error)
}
