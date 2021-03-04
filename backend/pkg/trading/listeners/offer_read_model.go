package listeners

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/commonpool/backend/pkg/eventbus"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/trading/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type OfferReadModel struct {
	ID      string
	Status  string
	Version int
}

type OfferItemReadModel struct {
	ID                     string `gorm:"not null;type:varchar(128);primaryKey"`
	OfferID                string `gorm:"not null;type:varchar(128)"`
	Version                int
	Type                   string `gorm:"not null;type:varchar(128)"`
	FromType               string `gorm:"type:varchar(128)"`
	FromID                 string `gorm:"type:varchar(128)"`
	ToType                 string `gorm:"type:varchar(128)"`
	ToID                   string `gorm:"type:varchar(128)"`
	ResourceID             string `gorm:"type:varchar(128)"`
	Amount                 time.Duration
	Duration               time.Duration
	AcceptedInbound        bool
	AcceptedInboundBy      string `gorm:"type:varchar(128)"`
	AcceptedInboundAt      time.Time
	AcceptedOutbound       bool
	AcceptedOutboundBy     string `gorm:"type:varchar(128)"`
	AcceptedOutbountAt     time.Time
	ServiceGiven           bool
	ServiceGivenBy         string `gorm:"type:varchar(128)"`
	ServiceGivenAt         time.Time
	ServiceReceived        bool
	ServiceReceivedBy      string `gorm:"type:varchar(128)"`
	ServiceReceivedAt      time.Time
	ResourceGiven          bool
	ResourceGivenBy        string `gorm:"type:varchar(128)"`
	ResourceGivenAt        time.Time
	ResourceTaken          bool
	ResourceTakenBy        string `gorm:"type:varchar(128)"`
	ResourceTakenAt        time.Time
	ResourceBorrowed       bool
	ResourceBorrowedBy     string `gorm:"type:varchar(128)"`
	ResourceBorrowedAt     time.Time
	ResourceLent           bool
	ResourceLentBy         string `gorm:"type:varchar(128)"`
	ResourceLentAt         time.Time
	BorrowedItemReturned   bool
	BorrowedItemReturnedBy string `gorm:"type:varchar(128)"`
	BorrowedItemReturnedAt time.Time
	LentItemReceived       bool
	LentItemReceivedBy     string `gorm:"type:varchar(128)"`
	LentItemReceivedAt     time.Time
}

type OfferReadModelHandler struct {
	catchUpFactory eventbus.CatchUpListenerFactory
	db             *gorm.DB
}

func NewOfferReadModelHandler(db *gorm.DB, catchUpFactory eventbus.CatchUpListenerFactory) *OfferReadModelHandler {
	var transactionHistoryHandler = &OfferReadModelHandler{
		catchUpFactory: catchUpFactory,
		db:             db,
	}
	return transactionHistoryHandler
}

func (h *OfferReadModelHandler) Start(ctx context.Context) error {

	if err := h.db.AutoMigrate(&OfferReadModel{}, &OfferItemReadModel{}); err != nil {
		return err
	}

	listener := h.catchUpFactory("transaction-history", time.Second*10)
	if err := listener.Initialize(ctx, "transaction-history", []string{
		string(domain.OfferSubmittedEvent),
	}); err != nil {
		return err
	}

	return listener.Listen(ctx, func(events []*eventstore.StreamEvent) error {
		for _, event := range events {
			if event.EventType == string(domain.OfferSubmittedEvent) {
				var evt domain.OfferSubmitted
				err := json.Unmarshal([]byte(event.Payload), &evt)
				if err != nil {
					return err
				}
				if err := h.HandleOfferSubmitted(&evt, event); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func getOptimisticLocking(db *gorm.DB, version int) *gorm.DB {
	return db.Clauses(
		clause.OnConflict{
			Where: clause.Where{
				Exprs: []clause.Expression{
					clause.Lt{
						Column: "version",
						Value:  version,
					},
				},
			},
			UpdateAll: true,
		})
}

func (h *OfferReadModelHandler) HandleOfferSubmitted(offerSubmitted *domain.OfferSubmitted, evt *eventstore.StreamEvent) error {

	return h.db.Transaction(func(tx *gorm.DB) error {

		getOptimisticLocking(tx, evt.SequenceNo).Create(&OfferReadModel{
			ID:      evt.StreamID,
			Status:  "pending",
			Version: evt.SequenceNo,
		})

		var offerItems []*OfferItemReadModel

		for _, offerItem := range offerSubmitted.OfferItems.Items {
			rm := &OfferItemReadModel{
				ID:      offerItem.Key().String(),
				OfferID: evt.StreamID,
				Type:    string(offerItem.Type()),
				Version: evt.SequenceNo,
			}
			if o, ok := offerItem.(*domain.ResourceTransferItem); ok {
				rm.ToType = string(o.To.Type)
				if o.To.GroupKey != nil {
					rm.ToID = o.To.GroupKey.String()
				} else {
					rm.ToID = o.To.UserKey.String()
				}
				rm.ResourceID = o.ResourceKey.String()
			} else if o, ok := offerItem.(*domain.ResourceBorrowItem); ok {
				rm.ToType = string(o.To.Type)
				if o.To.GroupKey != nil {
					rm.ToID = o.To.GroupKey.String()
				} else {
					rm.ToID = o.To.UserKey.String()
				}
				rm.ResourceID = o.ResourceKey.String()
				rm.Duration = o.Duration
			} else if o, ok := offerItem.(*domain.ServiceOfferItem); ok {
				rm.ToType = string(o.To.Type)
				if o.To.GroupKey != nil {
					rm.ToID = o.To.GroupKey.String()
				} else {
					rm.ToID = o.To.UserKey.String()
				}
				rm.FromType = string(o.From.Type)
				if o.From.GroupKey != nil {
					rm.FromID = o.From.GroupKey.String()
				} else {
					rm.FromID = o.From.UserKey.String()
				}
				rm.ResourceID = o.ResourceKey.String()
				rm.Duration = o.Duration
			} else if o, ok := offerItem.(*domain.CreditTransferItem); ok {
				rm.ToType = string(o.To.Type)
				if o.To.GroupKey != nil {
					rm.ToID = o.To.GroupKey.String()
				} else {
					rm.ToID = o.To.UserKey.String()
				}
				rm.FromType = string(o.From.Type)
				if o.From.GroupKey != nil {
					rm.FromID = o.From.GroupKey.String()
				} else {
					rm.FromID = o.From.UserKey.String()
				}
				rm.Amount = o.Amount
			}
			offerItems = append(offerItems, rm)
		}

		if err := getOptimisticLocking(tx, evt.SequenceNo).Create(offerItems).Error; err != nil {
			return err
		}

		return nil

	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
}
