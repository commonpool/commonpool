package service

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading"
	transaction2 "github.com/commonpool/backend/pkg/transaction"
	uuid "github.com/satori/go.uuid"
	"time"
)

type TransactionService struct {
	store transaction2.Store
}

func NewTransactionService(store transaction2.Store) *TransactionService {
	return &TransactionService{
		store: store,
	}
}

var _ transaction2.Service = &TransactionService{}

func (t TransactionService) UserSharedResourceWithGroup(groupKey keys.GroupKey, resourceKey keys.ResourceKey) (*transaction2.Entry, error) {

	entry := &transaction2.Entry{
		Key:         keys.NewTransactionEntryKey(uuid.NewV4()),
		Type:        transaction2.ResourceSharedWithGroup,
		GroupKey:    groupKey,
		ResourceKey: &resourceKey,
		Duration:    nil,
		Recipient:   nil,
		From:        nil,
		Timestamp:   time.Now().UTC(),
	}

	err := t.store.SaveEntry(entry)
	if err != nil {
		return nil, err
	}

	return t.store.GetEntry(entry.Key)

}

func (t TransactionService) UserRemovedResourceFromGroup(groupKey keys.GroupKey, resourceKey keys.ResourceKey) (*transaction2.Entry, error) {
	entry := &transaction2.Entry{
		Key:         keys.NewTransactionEntryKey(uuid.NewV4()),
		Type:        transaction2.ResourceRemovedFromGroup,
		GroupKey:    groupKey,
		ResourceKey: &resourceKey,
		Duration:    nil,
		Recipient:   nil,
		From:        nil,
		Timestamp:   time.Now().UTC(),
	}

	err := t.store.SaveEntry(entry)
	if err != nil {
		return nil, err
	}

	return t.store.GetEntry(entry.Key)
}

func (t TransactionService) ServiceWasProvided(groupKey keys.GroupKey, resourceKey keys.ResourceKey, duration time.Duration) (*transaction2.Entry, error) {
	entry := &transaction2.Entry{
		Key:         keys.NewTransactionEntryKey(uuid.NewV4()),
		Type:        transaction2.ServiceProvided,
		GroupKey:    groupKey,
		ResourceKey: &resourceKey,
		Duration:    &duration,
		Recipient:   nil,
		From:        nil,
		Timestamp:   time.Now().UTC(),
	}

	err := t.store.SaveEntry(entry)
	if err != nil {
		return nil, err
	}

	return t.store.GetEntry(entry.Key)
}

func (t TransactionService) ResourceWasBorrowed(groupKey keys.GroupKey, resourceKey keys.ResourceKey, recipient *trading.Target, expectedDuration time.Duration) (*transaction2.Entry, error) {
	entry := &transaction2.Entry{
		Key:         keys.NewTransactionEntryKey(uuid.NewV4()),
		Type:        transaction2.ResourceBorrowed,
		GroupKey:    groupKey,
		ResourceKey: &resourceKey,
		Duration:    &expectedDuration,
		Recipient:   recipient,
		From:        nil,
		Timestamp:   time.Now().UTC(),
	}

	err := t.store.SaveEntry(entry)
	if err != nil {
		return nil, err
	}

	return t.store.GetEntry(entry.Key)
}

func (t TransactionService) ResourceWasReturned(groupKey keys.GroupKey, resourceKey keys.ResourceKey, recipient *trading.Target, actualDuration time.Duration) (*transaction2.Entry, error) {
	entry := &transaction2.Entry{
		Key:         keys.NewTransactionEntryKey(uuid.NewV4()),
		Type:        transaction2.ResourceReturned,
		GroupKey:    groupKey,
		ResourceKey: &resourceKey,
		Duration:    &actualDuration,
		Recipient:   recipient,
		From:        nil,
		Timestamp:   time.Now().UTC(),
	}

	err := t.store.SaveEntry(entry)
	if err != nil {
		return nil, err
	}

	return t.store.GetEntry(entry.Key)
}

func (t TransactionService) ResourceWasTaken(groupKey keys.GroupKey, resourceKey keys.ResourceKey, recipient *trading.Target) (*transaction2.Entry, error) {
	entry := &transaction2.Entry{
		Key:         keys.NewTransactionEntryKey(uuid.NewV4()),
		Type:        transaction2.ResourceTaken,
		GroupKey:    groupKey,
		ResourceKey: &resourceKey,
		Duration:    nil,
		Recipient:   recipient,
		From:        nil,
		Timestamp:   time.Now().UTC(),
	}

	err := t.store.SaveEntry(entry)
	if err != nil {
		return nil, err
	}

	return t.store.GetEntry(entry.Key)
}

func (t TransactionService) TimeCreditsExchanged(groupKey keys.GroupKey, from *trading.Target, recipient *trading.Target, amount time.Duration) (*transaction2.Entry, error) {
	entry := &transaction2.Entry{
		Key:         keys.NewTransactionEntryKey(uuid.NewV4()),
		Type:        transaction2.TimeCreditsExchanged,
		GroupKey:    groupKey,
		ResourceKey: nil,
		Duration:    &amount,
		Recipient:   recipient,
		From:        from,
		Timestamp:   time.Now().UTC(),
	}

	err := t.store.SaveEntry(entry)
	if err != nil {
		return nil, err
	}

	return t.store.GetEntry(entry.Key)
}

func (t TransactionService) GetEntriesForGroupAndUsers(groupKey keys.GroupKey, userKeys *keys.UserKeys) (*transaction2.Entries, error) {
	return t.store.GetEntriesForGroupAndUsers(groupKey, userKeys)
}

func (t TransactionService) GetEntry(entryKey keys.TransactionEntryKey) (*transaction2.Entry, error) {
	return t.store.GetEntry(entryKey)
}
