package listeners

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/domain"
	"github.com/commonpool/backend/pkg/auth/readmodel"
	"github.com/commonpool/backend/pkg/eventbus"
	"github.com/commonpool/backend/pkg/eventsource"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type UserReadModelListener struct {
	db                     *gorm.DB
	catchUpListenerFactory eventbus.CatchUpListenerFactory
}

func NewUserReadModelListener(db *gorm.DB, catchUpListenerFactory eventbus.CatchUpListenerFactory) *UserReadModelListener {
	return &UserReadModelListener{
		db:                     db,
		catchUpListenerFactory: catchUpListenerFactory,
	}
}

func (l *UserReadModelListener) Start(ctx context.Context) error {

	if err := l.migrateDatabase(); err != nil {
		return err
	}

	catchUpListener := l.catchUpListenerFactory("user-read-model", time.Second*5)

	if err := catchUpListener.Initialize(ctx, "user-read-model", []string{
		domain.UserDiscoveredEvent,
		domain.UserInfoChangedEvent,
	}); err != nil {
		return err
	}

	return catchUpListener.Listen(ctx, l.handleEvents)
}

func (l *UserReadModelListener) handleEvents(events []eventsource.Event) error {
	for _, event := range events {
		if err := l.handleEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (l *UserReadModelListener) handleEvent(event eventsource.Event) error {
	switch e := event.(type) {
	case domain.UserDiscovered:
		return l.handleUserDiscovered(e)
	case domain.UserInfoChanged:
		return l.handleUserInfoChanged(e)
	}
	return nil
}

func (l *UserReadModelListener) handleUserDiscovered(e domain.UserDiscovered) error {
	return getOptimisticLocking(l.db, e.SequenceNo).
		Create(&readmodel.UserReadModel{
			UserKey:  e.AggregateID,
			Email:    e.UserInfo.Email,
			Username: e.UserInfo.Username,
			Version:  e.SequenceNo,
		}).Error
}

func (l *UserReadModelListener) handleUserInfoChanged(e domain.UserInfoChanged) error {
	return l.db.
		Model(&readmodel.UserReadModel{}).
		Where("user_key = ? and version < ?", e.AggregateID, e.SequenceNo).
		Updates(map[string]interface{}{
			"username": e.NewUserInfo.Username,
			"email":    e.NewUserInfo.Email,
		}).Error
}

func (l *UserReadModelListener) migrateDatabase() error {
	if err := l.db.AutoMigrate(&readmodel.UserReadModel{}); err != nil {
		return err
	}
	return nil
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
