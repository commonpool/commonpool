package transaction

import (
	"github.com/commonpool/backend/model"
	"time"
)

type TransactionType string

const (
	ResourceSharedWithGroup  TransactionType = "resource_shared_with_group"
	ResourceRemovedFromGroup TransactionType = "resource_removed_from_group"
	ServiceProvided          TransactionType = "service_provided"
	ResourceBorrowed         TransactionType = "resource_borrowed"
	ResourceReturned         TransactionType = "resource_returned"
	ResourceTaken            TransactionType = "resource_taken"
	TimeCreditsExchanged     TransactionType = "time_credits_exchanged"
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

type Entries struct {
	Items []*Entry
}

func NewEntries(items []*Entry) *Entries {
	return &Entries{
		Items: items,
	}
}

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

type Store interface {
	SaveEntry(entry *Entry) error
	GetEntry(entryKey model.TransactionEntryKey) (*Entry, error)
	GetEntriesForGroupAndUsers(groupKey model.GroupKey, userKeys *model.UserKeys) (*Entries, error)
}
