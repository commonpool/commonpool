package service

import (
	"github.com/commonpool/backend/model"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
	transaction2 "github.com/commonpool/backend/pkg/transaction"
	usermodel "github.com/commonpool/backend/pkg/user/model"
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

func (t TransactionService) UserSharedResourceWithGroup(groupKey groupmodel.GroupKey, resourceKey resourcemodel.ResourceKey) (*transaction2.Entry, error) {

	entry := &transaction2.Entry{
		Key:         model.NewTransactionEntryKey(uuid.NewV4()),
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

func (t TransactionService) UserRemovedResourceFromGroup(groupKey groupmodel.GroupKey, resourceKey resourcemodel.ResourceKey) (*transaction2.Entry, error) {
	entry := &transaction2.Entry{
		Key:         model.NewTransactionEntryKey(uuid.NewV4()),
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

func (t TransactionService) ServiceWasProvided(groupKey groupmodel.GroupKey, resourceKey resourcemodel.ResourceKey, duration time.Duration) (*transaction2.Entry, error) {
	entry := &transaction2.Entry{
		Key:         model.NewTransactionEntryKey(uuid.NewV4()),
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

func (t TransactionService) ResourceWasBorrowed(groupKey groupmodel.GroupKey, resourceKey resourcemodel.ResourceKey, recipient *resourcemodel.Target, expectedDuration time.Duration) (*transaction2.Entry, error) {
	entry := &transaction2.Entry{
		Key:         model.NewTransactionEntryKey(uuid.NewV4()),
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

func (t TransactionService) ResourceWasReturned(groupKey groupmodel.GroupKey, resourceKey resourcemodel.ResourceKey, recipient *resourcemodel.Target, actualDuration time.Duration) (*transaction2.Entry, error) {
	entry := &transaction2.Entry{
		Key:         model.NewTransactionEntryKey(uuid.NewV4()),
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

func (t TransactionService) ResourceWasTaken(groupKey groupmodel.GroupKey, resourceKey resourcemodel.ResourceKey, recipient *resourcemodel.Target) (*transaction2.Entry, error) {
	entry := &transaction2.Entry{
		Key:         model.NewTransactionEntryKey(uuid.NewV4()),
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

func (t TransactionService) TimeCreditsExchanged(groupKey groupmodel.GroupKey, from *resourcemodel.Target, recipient *resourcemodel.Target, amount time.Duration) (*transaction2.Entry, error) {
	entry := &transaction2.Entry{
		Key:         model.NewTransactionEntryKey(uuid.NewV4()),
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

func (t TransactionService) GetEntriesForGroupAndUsers(groupKey groupmodel.GroupKey, userKeys *usermodel.UserKeys) (*transaction2.Entries, error) {
	return t.store.GetEntriesForGroupAndUsers(groupKey, userKeys)
}

func (t TransactionService) GetEntry(entryKey model.TransactionEntryKey) (*transaction2.Entry, error) {
	return t.store.GetEntry(entryKey)
}
