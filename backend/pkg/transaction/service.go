package transaction

import (
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource"
	"time"
)

type Service interface {
	UserSharedResourceWithGroup(groupKey keys.GroupKey, resourceKey keys.ResourceKey) (*Entry, error)
	UserRemovedResourceFromGroup(groupKey keys.GroupKey, resourceKey keys.ResourceKey) (*Entry, error)
	ServiceWasProvided(groupKey keys.GroupKey, resourceKey keys.ResourceKey, duration time.Duration) (*Entry, error)
	ResourceWasBorrowed(groupKey keys.GroupKey, resourceKey keys.ResourceKey, recipient *resource.Target, expectedDuration time.Duration) (*Entry, error)
	ResourceWasReturned(groupKey keys.GroupKey, resourceKey keys.ResourceKey, recipient *resource.Target, actualDuration time.Duration) (*Entry, error)
	ResourceWasTaken(groupKey keys.GroupKey, resourceKey keys.ResourceKey, recipient *resource.Target) (*Entry, error)
	TimeCreditsExchanged(groupKey keys.GroupKey, from *resource.Target, recipient *resource.Target, amount time.Duration) (*Entry, error)
	GetEntriesForGroupAndUsers(groupKey keys.GroupKey, userKeys *keys.UserKeys) (*Entries, error)
	GetEntry(entryKey model.TransactionEntryKey) (*Entry, error)
}
