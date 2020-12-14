package transaction

import "github.com/commonpool/backend/model"

type Store interface {
	SaveEntry(entry *Entry) error
	GetEntry(entryKey model.TransactionEntryKey) (*Entry, error)
	GetEntriesForGroupAndUsers(groupKey model.GroupKey, userKeys *model.UserKeys) (*Entries, error)
}
