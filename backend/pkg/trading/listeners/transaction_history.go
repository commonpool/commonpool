package listeners

import (
	"context"
	"github.com/commonpool/backend/pkg/eventbus"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/trading/domain"
	"gorm.io/gorm"
	"time"
)

type CompletedOffer struct {
	ID string
}

type TransactionHistoryHandler struct {
	catchUpFactory eventbus.CatchUpListenerFactory
	db             *gorm.DB
}

func NewTransactionHistoryHandler(db *gorm.DB, catchUpFactory eventbus.CatchUpListenerFactory) *TransactionHistoryHandler {
	var transactionHistoryHandler = &TransactionHistoryHandler{
		catchUpFactory: catchUpFactory,
		db:             db,
	}
	return transactionHistoryHandler
}

func (h *TransactionHistoryHandler) Start(ctx context.Context) error {

	if err := h.db.AutoMigrate(&CompletedOffer{}); err != nil {
		return err
	}

	listener := h.catchUpFactory("transaction-history", time.Second*10)
	if err := listener.Initialize(ctx, "transaction-history", []string{
		domain.OfferCompletedEvent,
	}); err != nil {
		return err
	}
	return listener.Listen(ctx, func(ctx context.Context, events []eventsource.Event) error {
		for _, event := range events {
			if event.GetEventType() != domain.OfferCompletedEvent {
				continue
			}

			qry := h.db.Find(&CompletedOffer{}, "id = ?", event.GetAggregateID())
			if qry.Error != nil {
				return qry.Error
			}

			if qry.RowsAffected > 0 {
				return nil
			}

			newCompletedOffer := &CompletedOffer{
				ID: event.GetAggregateID(),
			}
			if err := h.db.Create(newCompletedOffer).Error; err != nil {
				return err
			}

		}
		return nil
	})
}
