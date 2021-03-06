package listeners

import (
	"context"
	"database/sql"
	"github.com/commonpool/backend/pkg/eventbus"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/trading/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type OfferReadModel struct {
	ID          string
	Status      string
	Version     int
	DeclinedAt  time.Time
	DeclinedBy  string
	SubmittedAt time.Time
	SubmittedBy string
	ApprovedAt  time.Time
	CompletedAt time.Time
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
	ApprovedInbound        bool
	ApprovedInboundBy      string `gorm:"type:varchar(128)"`
	ApprovedInboundAt      time.Time
	ApprovedOutbound       bool
	ApprovedOutboundBy     string `gorm:"type:varchar(128)"`
	ApprovedOutboundAt     time.Time
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

	listener := h.catchUpFactory("offer-read-model", time.Second*10)
	if err := listener.Initialize(ctx, "offer-read-model", []string{
		string(domain.OfferSubmittedEvent),
		string(domain.OfferApprovedEvent),
		string(domain.OfferDeclinedEvent),
		string(domain.OfferCompletedEvent),
		string(domain.OfferItemApprovedEvent),
		string(domain.ResourceTransferGivenNotifiedEvent),
		string(domain.ResourceTransferReceivedNotifiedEvent),
		string(domain.ServiceGivenNotifiedEvent),
		string(domain.ServiceReceivedNotifiedEvent),
		string(domain.ResourceBorrowedNotifiedEvent),
		string(domain.ResourceLentNotifiedEvent),
		string(domain.BorrowerReturnedResourceEvent),
		string(domain.LenderReceivedBackResourceEvent),
	}); err != nil {
		return err
	}

	return listener.Listen(ctx, func(events []eventsource.Event) error {
		for _, event := range events {

			switch e := event.(type) {
			case *domain.OfferSubmitted:
				if err := h.handleOfferSubmitted(e); err != nil {
					return err
				}
			case *domain.OfferApproved:
				if err := h.handleOfferApproved(e); err != nil {
					return err
				}
			case *domain.OfferDeclined:
				if err := h.handleOfferDeclined(e); err != nil {
					return err
				}
			case *domain.OfferCompleted:
				if err := h.handleOfferCompleted(e); err != nil {
					return err
				}
			case *domain.OfferItemApproved:
				updates := map[string]interface{}{
					"version": e.GetSequenceNo(),
				}
				if e.Direction == domain.Inbound {
					updates["approved_inbound"] = true
					updates["approved_inbound_by"] = e.ApprovedBy.String()
					updates["approved_inbound_at"] = e.GetEventTime()
				} else if e.Direction == domain.Outbound {
					updates["approved_outbound"] = true
					updates["approved_outbound_by"] = e.ApprovedBy.String()
					updates["approved_outbound_at"] = event.GetEventTime()
				}
				return h.updateOfferItem(e.OfferItemKey.String(), event.GetSequenceNo(), updates)
			case *domain.ResourceTransferGivenNotified:
				return h.updateOfferItem(e.OfferItemKey.String(), event.GetSequenceNo(), map[string]interface{}{
					"resource_given":    true,
					"resource_given_by": e.NotifiedBy,
					"resource_given_at": event.GetEventTime(),
				})
			case *domain.ResourceTransferReceivedNotified:
				return h.updateOfferItem(e.OfferItemKey.String(), event.GetSequenceNo(), map[string]interface{}{
					"resource_taken":    true,
					"resource_taken_by": e.NotifiedBy,
					"resource_taken_at": event.GetEventTime(),
				})
			case *domain.ServiceGivenNotified:
				return h.updateOfferItem(e.OfferItemKey.String(), event.GetSequenceNo(), map[string]interface{}{
					"service_given":    true,
					"service_given_by": e.ServiceGivenNotifiedPayload.NotifiedBy,
					"service_given_at": event.GetEventTime(),
				})
			case *domain.ServiceReceivedNotified:
				return h.updateOfferItem(e.OfferItemKey.String(), event.GetSequenceNo(), map[string]interface{}{
					"service_received":    true,
					"service_received_by": e.NotifiedBy,
					"service_received_at": event.GetEventTime(),
				})
			case *domain.ResourceBorrowedNotified:
				return h.updateOfferItem(e.OfferItemKey.String(), event.GetSequenceNo(), map[string]interface{}{
					"resource_borrowed":    true,
					"resource_borrowed_by": e.NotifiedBy,
					"resource_borrowed_at": event.GetEventTime(),
				})
			case *domain.ResourceLentNotified:
				return h.updateOfferItem(e.OfferItemKey.String(), event.GetSequenceNo(), map[string]interface{}{
					"resource_lent":    true,
					"resource_lent_by": e.NotifiedBy,
					"resource_lent_at": event.GetEventTime(),
				})
			case *domain.BorrowerReturnedResourceNotified:
				return h.updateOfferItem(e.OfferItemKey.String(), event.GetSequenceNo(), map[string]interface{}{
					"borrowed_item_returned":    true,
					"borrowed_item_returned_by": e.NotifiedBy,
					"borrowed_item_returned_at": event.GetEventTime(),
				})
			case *domain.LenderReceivedBackResourceNotified:
				return h.updateOfferItem(e.OfferItemKey.String(), event.GetSequenceNo(), map[string]interface{}{
					"lent_item_received":    true,
					"lent_item_received_by": e.NotifiedBy,
					"lent_item_received_at": event.GetEventTime(),
				})
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

func (h *OfferReadModelHandler) handleOfferCompleted(e *domain.OfferCompleted) error {
	return h.db.Transaction(func(tx *gorm.DB) error {
		return tx.
			Model(&OfferReadModel{}).
			Where("id = ? and version < ?", e.GetAggregateID(), e.GetSequenceNo()).
			Updates(map[string]interface{}{
				"status":       "completed",
				"version":      e.GetSequenceNo(),
				"completed_at": e.GetEventTime(),
			}).Error
	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
}

func (h *OfferReadModelHandler) handleOfferApproved(evt *domain.OfferApproved) error {
	return h.db.Transaction(func(tx *gorm.DB) error {
		return tx.
			Model(&OfferReadModel{}).
			Where("id = ? and version < ?", evt.GetAggregateID(), evt.GetSequenceNo()).
			Updates(map[string]interface{}{
				"status":      "approved",
				"version":     evt.GetSequenceNo(),
				"approved_at": evt.GetEventTime(),
			}).Error
	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
}

func (h *OfferReadModelHandler) handleOfferDeclined(e *domain.OfferDeclined) error {
	return h.db.Transaction(func(tx *gorm.DB) error {
		return tx.
			Model(&OfferReadModel{}).
			Where("id = ? and version < ?", e.GetAggregateID(), e.GetSequenceNo()).
			Updates(map[string]interface{}{
				"status":      "declined",
				"version":     e.GetSequenceNo(),
				"declined_by": e.OfferDeclinedPayload.DeclinedBy,
				"declined_at": e.GetEventTime(),
			}).Error
	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
}

func (h *OfferReadModelHandler) updateOfferItem(offerItemId string, expectedVersion int, updates map[string]interface{}) error {
	return h.db.Transaction(func(tx *gorm.DB) error {
		updates["version"] = expectedVersion
		return tx.
			Model(&OfferItemReadModel{}).
			Where("id = ? and version < ?", offerItemId, expectedVersion).
			Updates(updates).Error
	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
}

func (h *OfferReadModelHandler) handleOfferSubmitted(e *domain.OfferSubmitted) error {

	return h.db.Transaction(func(tx *gorm.DB) error {

		getOptimisticLocking(tx, e.GetSequenceNo()).Create(&OfferReadModel{
			ID:          e.GetAggregateID(),
			Status:      "pending",
			Version:     e.GetSequenceNo(),
			SubmittedBy: e.SubmittedBy.String(),
			SubmittedAt: e.GetEventTime(),
		})

		var offerItems []*OfferItemReadModel

		for _, offerItem := range e.OfferItems.Items {
			rm := &OfferItemReadModel{
				ID:      offerItem.GetKey().String(),
				OfferID: e.GetAggregateID(),
				Type:    string(offerItem.Type()),
				Version: e.GetSequenceNo(),
			}

			if resourceTransfer, ok := offerItem.AsResourceTransfer(); ok {
				rm.ToType = string(resourceTransfer.To.Type)
				rm.ToID = resourceTransfer.To.GetKeyAsString()
				rm.ResourceID = resourceTransfer.ResourceKey.String()
			}

			if provideService, ok := offerItem.AsProvideService(); ok {
				rm.ToType = string(provideService.To.Type)
				rm.ToID = provideService.To.GetKeyAsString()
				// TODO: rm.FromType = provideService.From.GetKeyAsString()
				// TODO: rm.FromID = provideService.From.GetKeyAsString()
				rm.ResourceID = provideService.ResourceKey.String()
				rm.Duration = provideService.Duration
			}

			if borrowResource, ok := offerItem.AsBorrowResource(); ok {
				rm.ToType = string(borrowResource.To.Type)
				rm.ToID = borrowResource.To.GetKeyAsString()
				// TODO: rm.FromID = provideServoce.From.GetKeyAsString()
				rm.ResourceID = borrowResource.ResourceKey.String()
				rm.Duration = borrowResource.Duration

			}

			if creditTransfer, ok := offerItem.AsCreditTransfer(); ok {
				rm.ToType = string(creditTransfer.To.Type)
				rm.ToID = creditTransfer.To.GetKeyAsString()
				rm.FromType = string(creditTransfer.From.Type)
				rm.FromID = creditTransfer.From.GetKeyAsString()
				rm.Amount = creditTransfer.Amount
			}

			offerItems = append(offerItems, rm)
		}

		if err := getOptimisticLocking(tx, e.GetSequenceNo()).Create(offerItems).Error; err != nil {
			return err
		}

		return nil

	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
}
