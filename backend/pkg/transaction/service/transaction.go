package service

import (
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/transaction"
	uuid "github.com/satori/go.uuid"
	"time"
)

type TransactionService struct {
	store transaction.Store
}

func NewTransactionService(store transaction.Store) *TransactionService {
	return &TransactionService{
		store: store,
	}
}

var _ transaction.Service = &TransactionService{}

func (t TransactionService) UserSharedResourceWithGroup(groupKey model.GroupKey, resourceKey model.ResourceKey) (*transaction.Entry, error) {

	entry := &transaction.Entry{
		Key:         model.NewTransactionEntryKey(uuid.NewV4()),
		Type:        transaction.ResourceSharedWithGroup,
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

func (t TransactionService) UserRemovedResourceFromGroup(groupKey model.GroupKey, resourceKey model.ResourceKey) (*transaction.Entry, error) {
	entry := &transaction.Entry{
		Key:         model.NewTransactionEntryKey(uuid.NewV4()),
		Type:        transaction.ResourceRemovedFromGroup,
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

func (t TransactionService) ServiceWasProvided(groupKey model.GroupKey, resourceKey model.ResourceKey, duration time.Duration) (*transaction.Entry, error) {
	entry := &transaction.Entry{
		Key:         model.NewTransactionEntryKey(uuid.NewV4()),
		Type:        transaction.ResourceRemovedFromGroup,
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

func (t TransactionService) ResourceWasBorrowed(groupKey model.GroupKey, resourceKey model.ResourceKey, recipient *model.Target, expectedDuration time.Duration) (*transaction.Entry, error) {
	entry := &transaction.Entry{
		Key:         model.NewTransactionEntryKey(uuid.NewV4()),
		Type:        transaction.ResourceBorrowed,
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

func (t TransactionService) ResourceWasReturned(groupKey model.GroupKey, resourceKey model.ResourceKey, recipient *model.Target, actualDuration time.Duration) (*transaction.Entry, error) {
	entry := &transaction.Entry{
		Key:         model.NewTransactionEntryKey(uuid.NewV4()),
		Type:        transaction.ResourceReturned,
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

func (t TransactionService) ResourceWasTaken(groupKey model.GroupKey, resourceKey model.ResourceKey, recipient *model.Target) (*transaction.Entry, error) {
	entry := &transaction.Entry{
		Key:         model.NewTransactionEntryKey(uuid.NewV4()),
		Type:        transaction.ResourceTaken,
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

func (t TransactionService) TimeCreditsExchanged(groupKey model.GroupKey, from *model.Target, recipient *model.Target, amount time.Duration) (*transaction.Entry, error) {
	entry := &transaction.Entry{
		Key:         model.NewTransactionEntryKey(uuid.NewV4()),
		Type:        transaction.TimeCreditsExchanged,
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

func (t TransactionService) GetEntriesForGroupAndUsers(groupKey model.GroupKey, userKeys *model.UserKeys) (*transaction.Entries, error) {
	return t.store.GetEntriesForGroupAndUsers(groupKey, userKeys)
}

func (t TransactionService) GetEntry(entryKey model.TransactionEntryKey) (*transaction.Entry, error) {
	return t.store.GetEntry(entryKey)
}
