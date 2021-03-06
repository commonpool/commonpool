package transaction

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
	"time"
)

type Service interface {
	UserSharedResourceWithGroup(groupKey keys.GroupKey, resourceKey keys.ResourceKey) (*Entry, error)
	UserRemovedResourceFromGroup(groupKey keys.GroupKey, resourceKey keys.ResourceKey) (*Entry, error)
	ServiceWasProvided(groupKey keys.GroupKey, resourceKey keys.ResourceKey, duration time.Duration) (*Entry, error)
	ResourceWasBorrowed(groupKey keys.GroupKey, resourceKey keys.ResourceKey, recipient *domain.Target, expectedDuration time.Duration) (*Entry, error)
	ResourceWasReturned(groupKey keys.GroupKey, resourceKey keys.ResourceKey, recipient *domain.Target, actualDuration time.Duration) (*Entry, error)
	ResourceWasTaken(groupKey keys.GroupKey, resourceKey keys.ResourceKey, recipient *domain.Target) (*Entry, error)
	TimeCreditsExchanged(groupKey keys.GroupKey, from *domain.Target, recipient *domain.Target, amount time.Duration) (*Entry, error)
	GetEntriesForGroupAndUsers(groupKey keys.GroupKey, userKeys *keys.UserKeys) (*Entries, error)
	GetEntry(entryKey keys.TransactionEntryKey) (*Entry, error)
}
