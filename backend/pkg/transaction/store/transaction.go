package store

import (
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource"
	"github.com/commonpool/backend/pkg/transaction"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type TransactionStore struct {
	db *gorm.DB
}

func NewTransactionStore(db *gorm.DB) *TransactionStore {
	return &TransactionStore{
		db: db,
	}
}

var _ transaction.Store = &TransactionStore{}

func (t TransactionStore) SaveEntry(entry *transaction.Entry) error {

	var resourceID *uuid.UUID = nil
	if entry.ResourceKey != nil {
		resourceKeyVal := entry.ResourceKey.ID
		resourceID = &resourceKeyVal
	}

	var recipientType *resource.TargetType
	if entry.Recipient != nil {
		recipientTypeVal := entry.Recipient.Type
		recipientType = &recipientTypeVal
	}

	var fromType *resource.TargetType
	if entry.From != nil {
		fromTypeVal := entry.From.Type
		fromType = &fromTypeVal
	}

	var recipientID *string
	if entry.Recipient != nil {
		if entry.Recipient.IsForGroup() {
			recipientIDVal := entry.Recipient.GroupKey.String()
			recipientID = &recipientIDVal
		} else if entry.Recipient.IsForUser() {
			recipientIDVal := entry.Recipient.UserKey.String()
			recipientID = &recipientIDVal
		}
	}

	var fromID *string
	if entry.From != nil {
		if entry.From.IsForGroup() {
			fromIdVal := entry.From.GroupKey.String()
			fromID = &fromIdVal
		} else if entry.From.IsForUser() {
			fromIdVal := entry.From.UserKey.String()
			fromID = &fromIdVal
		}
	}

	dbEntry := TransactionEntry{
		ID:            entry.Key.ID,
		Type:          entry.Type,
		GroupID:       entry.GroupKey.ID,
		ResourceID:    resourceID,
		Duration:      entry.Duration,
		RecipientType: recipientType,
		RecipientID:   recipientID,
		FromType:      fromType,
		FromID:        fromID,
		Timestamp:     entry.Timestamp,
	}

	err := t.db.Create(&dbEntry).Error
	if err != nil {
		return err
	}

	return nil
}

func (t TransactionStore) GetEntry(transactionKey model.TransactionEntryKey) (*transaction.Entry, error) {
	var transactionEntry TransactionEntry
	err := t.db.Model(TransactionEntry{}).First(&transactionEntry, "id = ?", transactionKey.String()).Error
	if err != nil {
		return nil, err
	}
	return mapDbTransactionEntry(&transactionEntry)
}

func (t TransactionStore) GetEntriesForGroupAndUsers(groupKey keys.GroupKey, userKeys *keys.UserKeys) (*transaction.Entries, error) {

	sql := "(group_id = ? OR recipient_id = ? OR from_id = ?)"
	var params []interface{}
	params = append(params, groupKey.String())
	params = append(params, groupKey.String())
	params = append(params, groupKey.String())

	if !userKeys.IsEmpty() {
		sql = sql + " OR ( recipient_id is not null and recipient_id in ("
		for _, userKey := range userKeys.Items {
			sql = sql + "?"
			params = append(params, userKey.String())
		}
		sql = sql + ") or (from_id is not null and from_id in ("
		for _, userKey := range userKeys.Items {
			sql = sql + "?"
			params = append(params, userKey.String())
		}
		sql = sql + ")"
	}

	sql = sql + " "

	var dbEntries []*TransactionEntry
	t.db.Where(sql, params...).Find(&dbEntries)

	var entries []*transaction.Entry
	for _, dbEntry := range dbEntries {
		entry, err := mapDbTransactionEntry(dbEntry)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return transaction.NewEntries(entries), nil

}

type TransactionEntry struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key"`
	Type          transaction.Type
	GroupID       uuid.UUID
	ResourceID    *uuid.UUID
	Duration      *time.Duration
	RecipientType *resource.TargetType
	RecipientID   *string
	FromType      *resource.TargetType
	FromID        *string
	Timestamp     time.Time
}

func mapDbTransactionEntry(dbTransactionEntry *TransactionEntry) (*transaction.Entry, error) {

	var resourceKey *keys.ResourceKey = nil
	if dbTransactionEntry.ResourceID != nil {
		resourceKeyVal := keys.NewResourceKey(*dbTransactionEntry.ResourceID)
		resourceKey = &resourceKeyVal
	}

	recipient, err := mapTarget(dbTransactionEntry.RecipientType, dbTransactionEntry.RecipientID)
	if err != nil {
		return nil, err
	}

	from, err := mapTarget(dbTransactionEntry.FromType, dbTransactionEntry.FromID)
	if err != nil {
		return nil, err
	}

	return &transaction.Entry{
		Key:         model.NewTransactionEntryKey(dbTransactionEntry.ID),
		Type:        dbTransactionEntry.Type,
		GroupKey:    keys.NewGroupKey(dbTransactionEntry.GroupID),
		ResourceKey: resourceKey,
		Duration:    dbTransactionEntry.Duration,
		Recipient:   recipient,
		From:        from,
		Timestamp:   dbTransactionEntry.Timestamp,
	}, nil

}

func mapTarget(targetType *resource.TargetType, targetId *string) (*resource.Target, error) {
	var target *resource.Target
	if targetType != nil && targetId != nil {
		if targetType.IsUser() {
			target = resource.NewUserTarget(keys.NewUserKey(*targetId))
		} else if targetType.IsGroup() {
			groupKey, err := keys.ParseGroupKey(*targetId)
			if err != nil {
				return nil, err
			}
			target = resource.NewGroupTarget(groupKey)
		}
	}
	return target, nil
}
