package transaction

import (
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/group"
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
	"time"
)

type Service interface {
	UserSharedResourceWithGroup(groupKey group.GroupKey, resourceKey resourcemodel.ResourceKey) (*Entry, error)
	UserRemovedResourceFromGroup(groupKey group.GroupKey, resourceKey resourcemodel.ResourceKey) (*Entry, error)
	ServiceWasProvided(groupKey group.GroupKey, resourceKey resourcemodel.ResourceKey, duration time.Duration) (*Entry, error)
	ResourceWasBorrowed(groupKey group.GroupKey, resourceKey resourcemodel.ResourceKey, recipient *resourcemodel.Target, expectedDuration time.Duration) (*Entry, error)
	ResourceWasReturned(groupKey group.GroupKey, resourceKey resourcemodel.ResourceKey, recipient *resourcemodel.Target, actualDuration time.Duration) (*Entry, error)
	ResourceWasTaken(groupKey group.GroupKey, resourceKey resourcemodel.ResourceKey, recipient *resourcemodel.Target) (*Entry, error)
	TimeCreditsExchanged(groupKey group.GroupKey, from *resourcemodel.Target, recipient *resourcemodel.Target, amount time.Duration) (*Entry, error)
	GetEntriesForGroupAndUsers(groupKey group.GroupKey, userKeys *usermodel.UserKeys) (*Entries, error)
	GetEntry(entryKey model.TransactionEntryKey) (*Entry, error)
}
