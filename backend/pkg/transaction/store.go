package transaction

import (
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/group"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
)

type Store interface {
	SaveEntry(entry *Entry) error
	GetEntry(entryKey model.TransactionEntryKey) (*Entry, error)
	GetEntriesForGroupAndUsers(groupKey group.GroupKey, userKeys *usermodel.UserKeys) (*Entries, error)
}
