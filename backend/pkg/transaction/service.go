package transaction

import (
	"github.com/commonpool/backend/model"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
	"time"
)

type Service interface {
	UserSharedResourceWithGroup(groupKey groupmodel.GroupKey, resourceKey resourcemodel.ResourceKey) (*Entry, error)
	UserRemovedResourceFromGroup(groupKey groupmodel.GroupKey, resourceKey resourcemodel.ResourceKey) (*Entry, error)
	ServiceWasProvided(groupKey groupmodel.GroupKey, resourceKey resourcemodel.ResourceKey, duration time.Duration) (*Entry, error)
	ResourceWasBorrowed(groupKey groupmodel.GroupKey, resourceKey resourcemodel.ResourceKey, recipient *resourcemodel.Target, expectedDuration time.Duration) (*Entry, error)
	ResourceWasReturned(groupKey groupmodel.GroupKey, resourceKey resourcemodel.ResourceKey, recipient *resourcemodel.Target, actualDuration time.Duration) (*Entry, error)
	ResourceWasTaken(groupKey groupmodel.GroupKey, resourceKey resourcemodel.ResourceKey, recipient *resourcemodel.Target) (*Entry, error)
	TimeCreditsExchanged(groupKey groupmodel.GroupKey, from *resourcemodel.Target, recipient *resourcemodel.Target, amount time.Duration) (*Entry, error)
	GetEntriesForGroupAndUsers(groupKey groupmodel.GroupKey, userKeys *usermodel.UserKeys) (*Entries, error)
	GetEntry(entryKey model.TransactionEntryKey) (*Entry, error)
}
