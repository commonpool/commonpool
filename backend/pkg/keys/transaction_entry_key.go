package keys

import "github.com/satori/go.uuid"

type TransactionEntryKey struct {
	ID uuid.UUID
}

func (k TransactionEntryKey) String() string {
	return k.ID.String()
}

func NewTransactionEntryKey(uid uuid.UUID) TransactionEntryKey {
	return TransactionEntryKey{
		ID: uid,
	}
}
