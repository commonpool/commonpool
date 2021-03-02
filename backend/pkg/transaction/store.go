package transaction

import (
	"github.com/commonpool/backend/pkg/keys"
)

type Store interface {
	SaveEntry(entry *Entry) error
	GetEntry(entryKey keys.TransactionEntryKey) (*Entry, error)
	GetEntriesForGroupAndUsers(groupKey keys.GroupKey, userKeys *keys.UserKeys) (*Entries, error)
}
