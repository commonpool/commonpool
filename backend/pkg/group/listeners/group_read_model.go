package listeners

import (
	"context"
	"database/sql"
	"github.com/commonpool/backend/pkg/eventbus"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/group/readmodels"
	"github.com/commonpool/backend/pkg/keys"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type GroupReadModelListener struct {
	catchUpListenerFactory eventbus.CatchUpListenerFactory
	db                     *gorm.DB
}

func NewGroupReadModelListener(catchUpListenerFactory eventbus.CatchUpListenerFactory, db *gorm.DB) *GroupReadModelListener {
	return &GroupReadModelListener{
		catchUpListenerFactory: catchUpListenerFactory,
		db:                     db,
	}
}

func (l *GroupReadModelListener) Start(ctx context.Context) error {

	if err := l.migrateDatabase(); err != nil {
		return err
	}

	factory := l.catchUpListenerFactory("locks.readmodels.group", time.Second*5)
	if err := factory.Initialize(ctx, "listeners.group-readmodel", []string{
		domain.GroupCreatedEvent,
		domain.GroupInfoChangedEvent,
		domain.MembershipStatusChangedEvent,
	}); err != nil {
		return err
	}
	return factory.Listen(ctx, l.applyEvents)
}

func (l *GroupReadModelListener) applyEvents(events []eventsource.Event) error {
	for _, event := range events {
		if err := l.applyEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (l *GroupReadModelListener) migrateDatabase() error {
	return l.db.AutoMigrate(&readmodels.GroupReadModel{}, &readmodels.MembershipReadModel{})
}

func (l *GroupReadModelListener) applyEvent(event eventsource.Event) error {
	switch e := event.(type) {
	case domain.GroupCreated:
		return l.applyGroupCreatedEvent(e)
	case domain.GroupInfoChanged:
		return l.applyGroupInfoChangedEvent(e)
	case domain.MembershipStatusChanged:
		return l.applyMembershipStatusChangedEvent(e)
	}
	return nil
}

func (l *GroupReadModelListener) applyMembershipStatusChangedEvent(e domain.MembershipStatusChanged) error {
	if e.IsNewMembership {

		updates := map[string]interface{}{
			"version":            e.SequenceNo,
			"group_key":          e.AggregateID,
			"user_key":           e.MemberKey.String(),
			"user_confirmed":     false,
			"user_confirmed_at":  nil,
			"group_confirmed":    false,
			"group_confirmed_by": nil,
			"group_confirmed_at": nil,
			"group_name":         e.GroupName,
		}

		applyMembershipChanged(updates, e)

		err := l.db.Clauses(
			clause.OnConflict{
				Where: clause.Where{
					Exprs: []clause.Expression{
						clause.Lt{
							Column: "version",
							Value:  e.SequenceNo,
						},
					},
				},
				UpdateAll: true,
			}).Model(&readmodels.MembershipReadModel{}).Create(updates).Error

		if err != nil {
			return err
		}

	} else if e.IsCanceledMembership {
		err := l.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.
				Where("group_key = ? and user_key = ?", e.AggregateID, e.MemberKey.String()).
				Delete(&readmodels.MembershipReadModel{}).
				Error; err != nil {
				return err
			}
			return nil
		}, &sql.TxOptions{
			Isolation: sql.LevelSerializable,
		})
		if err != nil {
			return err
		}
	} else {

		updates := map[string]interface{}{
			"version": e.SequenceNo,
		}

		applyMembershipChanged(updates, e)

		if err := l.db.Model(&readmodels.MembershipReadModel{}).
			Where("group_key = ? and user_key = ? and version < ?", e.AggregateID, e.MemberKey.String(), e.SequenceNo).
			Updates(updates).Error; err != nil {
			return err
		}
	}
	return nil
}

func (l *GroupReadModelListener) applyGroupInfoChangedEvent(e domain.GroupInfoChanged) error {

	return l.db.Transaction(func(tx *gorm.DB) error {

		err := tx.
			Model(&readmodels.GroupReadModel{}).
			Where("group_key = ? and version < ?", e.AggregateID, e.SequenceNo).
			Updates(map[string]interface{}{
				"name":        e.NewGroupInfo.Name,
				"description": e.NewGroupInfo.Description,
				"version":     e.SequenceNo,
			}).Error
		if err != nil {
			return err
		}

		if e.OldGroupInfo.Name != e.NewGroupInfo.Name {
			qry := tx.
				Model(&readmodels.MembershipReadModel{}).
				Where("group_key = ? and version < ?", e.AggregateID, e.SequenceNo).
				Updates(map[string]interface{}{
					"group_name": e.NewGroupInfo.Name,
					"version":    e.SequenceNo,
				})
			err := qry.Error
			if err != nil {
				return err
			}
		}

		return nil

	}, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

}

func (l *GroupReadModelListener) applyGroupCreatedEvent(e domain.GroupCreated) error {
	groupKey, err := keys.ParseGroupKey(e.AggregateID)
	if err != nil {
		return err
	}

	return l.db.Clauses(
		clause.OnConflict{
			Where: clause.Where{
				Exprs: []clause.Expression{
					clause.Lt{
						Column: "version",
						Value:  e.SequenceNo,
					},
				},
			},
			UpdateAll: true,
		}).Create(&readmodels.GroupReadModel{
		Version:     e.SequenceNo,
		GroupKey:    groupKey,
		Name:        e.GroupInfo.Name,
		Description: e.GroupInfo.Description,
		CreatedBy:   e.CreatedBy.String(),
		CreatedAt:   e.EventTime,
	}).Error
}

func applyMembershipChanged(values map[string]interface{}, e domain.MembershipStatusChanged) {

	if e.OldStatus == nil && *e.NewStatus == domain.ApprovedMembershipStatus {
		values["group_confirmed"] = true
		values["group_confirmed_by"] = e.ChangedBy.String()
		values["group_confirmed_at"] = e.EventTime
		values["user_confirmed"] = true
		values["user_confirmed_at"] = e.EventTime
	} else if (e.OldStatus == nil && *e.NewStatus == domain.PendingGroupMembershipStatus) ||
		(*e.NewStatus == domain.ApprovedMembershipStatus && *e.OldStatus == domain.PendingUserMembershipStatus) {
		values["user_confirmed"] = true
		values["user_confirmed_at"] = e.EventTime
	} else if (e.OldStatus == nil && *e.NewStatus == domain.PendingUserMembershipStatus) ||
		(*e.NewStatus == domain.ApprovedMembershipStatus && *e.OldStatus == domain.PendingGroupMembershipStatus) {
		values["group_confirmed"] = true
		values["group_confirmed_by"] = e.ChangedBy.String()
		values["group_confirmed_at"] = e.EventTime
	}

	if e.IsNewMembership || e.OldPermissions != e.NewPermissions {
		values["is_admin"] = e.NewPermissions.IsAdmin()
		values["is_owner"] = e.NewPermissions.IsOwner()
		values["is_member"] = e.NewPermissions.IsMember()
	}

	values["status"] = e.NewStatus

}
