package transaction

import (
	"github.com/commonpool/backend/model"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	usermodel "github.com/commonpool/backend/pkg/user/model"
)

type Store interface {
	SaveEntry(entry *Entry) error
	GetEntry(entryKey model.TransactionEntryKey) (*Entry, error)
	GetEntriesForGroupAndUsers(groupKey groupmodel.GroupKey, userKeys *usermodel.UserKeys) (*Entries, error)
}
