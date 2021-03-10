package listeners

import (
	"context"
	"database/sql"
	userdomain "github.com/commonpool/backend/pkg/auth/domain"
	"github.com/commonpool/backend/pkg/eventbus"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/group/readmodels"
	"github.com/commonpool/backend/pkg/keys"
	"golang.org/x/sync/errgroup"
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
	watchedEvents := domain.AllEvents
	watchedEvents = append(watchedEvents, userdomain.UserDiscoveredEvent)
	watchedEvents = append(watchedEvents, userdomain.UserInfoChangedEvent)
	if err := factory.Initialize(ctx, "listeners.group-readmodel", watchedEvents); err != nil {
		return err
	}
	return factory.Listen(ctx, l.applyEvents)
}

func (l *GroupReadModelListener) applyEvents(ctx context.Context, events []eventsource.Event) error {
	for _, event := range events {
		if err := l.applyEvent(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

func (l *GroupReadModelListener) migrateDatabase() error {
	return l.db.AutoMigrate(&readmodels.GroupReadModel{}, &readmodels.MembershipReadModel{}, &readmodels.DBGroupUserReadModel{})
}

func (l *GroupReadModelListener) applyEvent(ctx context.Context, event eventsource.Event) error {
	switch e := event.(type) {
	case domain.GroupCreated:
		return l.applyGroupCreatedEvent(e)
	case domain.GroupInfoChanged:
		return l.applyGroupInfoChangedEvent(ctx, e)
	case domain.MembershipStatusChanged:
		return l.applyMembershipStatusChangedEvent(e)
	case userdomain.UserDiscovered:
		return l.handleUserDiscovered(ctx, e)
	case userdomain.UserInfoChanged:
		return l.handleUserInfoChanged(ctx, e)
	}
	return nil
}

func (l *GroupReadModelListener) applyMembershipStatusChangedEvent(e domain.MembershipStatusChanged) error {

	var memberName string
	var memberVersion = -1
	var changedByName string
	var changedByVersion = -1

	if e.IsNewMembership || !e.IsCanceledMembership {
		var users []*readmodels.DBGroupUserReadModel
		qry := l.db.Model(&readmodels.DBGroupUserReadModel{}).Where("user_key in (?, ?)", e.MemberKey, e.ChangedBy).Find(&users)
		err := qry.Error
		if err != nil {
			return err
		}
		for _, user := range users {
			if user.UserKey == e.MemberKey {
				memberName = user.Name
				memberVersion = user.Version
			}
			if user.UserKey == e.ChangedBy {
				changedByName = user.Name
				changedByVersion = user.Version
			}
		}
	}

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
			"user_name":          memberName,
			"user_version":       memberVersion,
			"created_at":         e.EventTime,
			"created_by":         e.ChangedBy,
			"created_by_name":    changedByName,
			"created_by_version": changedByVersion,
		}
		populateMembershipUpdates(updates, e)
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
			}).Model(&readmodels.MembershipReadModel{}).Create(updates).Error
	} else if e.IsCanceledMembership {
		return l.db.Transaction(func(tx *gorm.DB) error {
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
	} else {
		updates := map[string]interface{}{
			"version": e.SequenceNo,
		}
		populateMembershipUpdates(updates, e)
		return l.db.Model(&readmodels.MembershipReadModel{}).
			Where("group_key = ? and user_key = ? and version < ?", e.AggregateID, e.MemberKey.String(), e.SequenceNo).
			Updates(updates).Error
	}
}

func (l *GroupReadModelListener) applyGroupInfoChangedEvent(ctx context.Context, e domain.GroupInfoChanged) error {
	return l.db.Transaction(func(tx *gorm.DB) error {
		g, _ := errgroup.WithContext(ctx)
		g.Go(func() error {
			return tx.
				Model(&readmodels.GroupReadModel{}).
				Where("group_key = ? and version < ?", e.AggregateID, e.SequenceNo).
				Updates(map[string]interface{}{
					"name":        e.NewGroupInfo.Name,
					"description": e.NewGroupInfo.Description,
					"version":     e.SequenceNo,
				}).Error
		})
		g.Go(func() error {
			if e.OldGroupInfo.Name != e.NewGroupInfo.Name {
				return tx.
					Model(&readmodels.MembershipReadModel{}).
					Where("group_key = ? and version < ?", e.AggregateID, e.SequenceNo).
					Updates(map[string]interface{}{
						"group_name": e.NewGroupInfo.Name,
						"version":    e.SequenceNo,
					}).Error
			}
			return nil
		})
		return g.Wait()
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

func (l *GroupReadModelListener) handleUserDiscovered(ctx context.Context, e userdomain.UserDiscovered) error {
	userKey := keys.NewUserKey(e.AggregateID)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return getOptimisticLocking(l.db, e.SequenceNo).Model(&readmodels.DBGroupUserReadModel{}).Create(&readmodels.DBGroupUserReadModel{
			UserKey: userKey,
			Name:    e.UserInfo.Username,
			Version: e.SequenceNo,
		}).Error
	})
	g.Go(func() error {
		return l.db.Model(&readmodels.MembershipReadModel{}).Where("user_key = ? and user_version < ?", userKey, e.SequenceNo).
			Updates(map[string]interface{}{
				"user_name":    e.UserInfo.Username,
				"user_version": e.SequenceNo,
			}).Error
	})
	g.Go(func() error {
		return l.db.Model(&readmodels.MembershipReadModel{}).Where("created_by = ? and created_by_version < ?", userKey, e.SequenceNo).
			Updates(map[string]interface{}{
				"created_by_name":    e.UserInfo.Username,
				"created_by_version": e.SequenceNo,
			}).Error
	})
	return g.Wait()
}

func (l *GroupReadModelListener) handleUserInfoChanged(ctx context.Context, e userdomain.UserInfoChanged) error {
	userKey := keys.NewUserKey(e.AggregateID)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return getOptimisticLocking(l.db, e.SequenceNo).Model(&readmodels.DBGroupUserReadModel{}).Create(&readmodels.DBGroupUserReadModel{
			UserKey: userKey,
			Name:    e.NewUserInfo.Username,
			Version: e.SequenceNo,
		}).Error
	})
	g.Go(func() error {
		return l.db.Where("user_key = ? and user_version < ?", userKey, e.SequenceNo).Model(&readmodels.MembershipReadModel{}).Updates(map[string]interface{}{
			"user_name":    e.NewUserInfo.Username,
			"user_version": e.SequenceNo,
		}).Error
	})
	g.Go(func() error {
		return l.db.Model(&readmodels.MembershipReadModel{}).Where("created_by = ? and created_by_version < ?", userKey, e.SequenceNo).
			Updates(map[string]interface{}{
				"created_by_name":    e.NewUserInfo.Username,
				"created_by_version": e.SequenceNo,
			}).Error
	})
	return g.Wait()
}

func populateMembershipUpdates(values map[string]interface{}, e domain.MembershipStatusChanged) {

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
