package listeners

import (
	"context"
	userdomain "github.com/commonpool/backend/pkg/auth/domain"
	"github.com/commonpool/backend/pkg/eventbus"
	"github.com/commonpool/backend/pkg/eventsource"
	groupdomain "github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/domain"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	"github.com/labstack/gommon/log"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ResourceReadModelHandler struct {
	db                     *gorm.DB
	catchUpListenerFactory eventbus.CatchUpListenerFactory
}

func NewResourceReadModelHandler(db *gorm.DB, catchUpListenerFactory eventbus.CatchUpListenerFactory) *ResourceReadModelHandler {
	return &ResourceReadModelHandler{db: db, catchUpListenerFactory: catchUpListenerFactory}
}

func (l *ResourceReadModelHandler) Start(ctx context.Context) error {
	err2 := l.migrateDatabase()
	if err2 != nil {
		return err2
	}
	listener := l.catchUpListenerFactory("readmodels.resource", time.Second*5)
	eventTypes := domain.AllEventTypes
	eventTypes = append(eventTypes, userdomain.UserDiscoveredEvent)
	eventTypes = append(eventTypes, userdomain.UserInfoChangedEvent)
	eventTypes = append(eventTypes, groupdomain.GroupCreatedEvent)
	eventTypes = append(eventTypes, groupdomain.GroupInfoChangedEvent)
	err := listener.Initialize(ctx, "readmodels.resource", eventTypes)
	if err != nil {
		return err
	}
	return listener.Listen(ctx, l.handleEvents)
}

func (l *ResourceReadModelHandler) migrateDatabase() error {
	if err := l.db.AutoMigrate(
		&readmodel.DbResourceReadModel{},
		&readmodel.ResourceSharingReadModel{},
		&readmodel.ResourceUserNameReadModel{},
		&readmodel.ResourceGroupNameReadModel{},
		&readmodel.ResourceEvaluationReadModel{},
	); err != nil {
		return err
	}
	return nil
}

func (l *ResourceReadModelHandler) handleEvents(ctx context.Context, events []eventsource.Event) error {
	for _, event := range events {
		if err := l.handleEvent(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

func (l *ResourceReadModelHandler) handleEvent(ctx context.Context, event eventsource.Event) error {
	switch e := event.(type) {
	case groupdomain.GroupCreated:
		return l.handleGroupCreated(ctx, e)
	case groupdomain.GroupInfoChanged:
		return l.handleGroupInfoChanged(ctx, e)
	case userdomain.UserDiscovered:
		return l.handleUserDiscovered(ctx, e)
	case userdomain.UserInfoChanged:
		return l.handleUserInfoChanged(ctx, e)
	case domain.ResourceRegistered:
		return l.handleResourceRegistered(e)
	case domain.ResourceInfoChanged:
		return l.handleResourceInfoChanged(e)
	case domain.ResourceGroupSharingChanged:
		return l.handleResourceGroupSharingChanged(ctx, e)
	case domain.ResourceDeleted:
		return l.handleResourceDeleted(ctx, e)
	case domain.ResourceEvaluated:
		return l.handleResourceEvaluated(ctx, e)
	default:
		log.Warnf("unhandled event type: %s", event.GetEventType())
	}
	return nil
}

func (l *ResourceReadModelHandler) handleGroupCreated(ctx context.Context, e groupdomain.GroupCreated) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return getOptimisticLocking(l.db, e.SequenceNo).
			Model(&readmodel.ResourceGroupNameReadModel{}).
			Create(readmodel.ResourceGroupNameReadModel{
				GroupKey:  e.AggregateID,
				GroupName: e.GroupCreatedPayload.GroupInfo.Name,
				Version:   e.SequenceNo,
			}).Error
	})
	g.Go(func() error {
		return l.doGroupUpdate(e.AggregateID, e.GroupInfo.Name, e.SequenceNo)
	})
	return g.Wait()
}

func (l *ResourceReadModelHandler) handleGroupInfoChanged(ctx context.Context, e groupdomain.GroupInfoChanged) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return l.db.
			Model(&readmodel.ResourceGroupNameReadModel{}).
			Where("group_key = ? and version <= ?", e.AggregateID, e.SequenceNo).
			Updates(map[string]interface{}{
				"group_name": e.NewGroupInfo.Name,
				"version":    e.SequenceNo,
			}).Error
	})
	g.Go(func() error {
		return l.doGroupUpdate(e.AggregateID, e.NewGroupInfo.Name, e.SequenceNo)
	})
	return g.Wait()
}

func (l *ResourceReadModelHandler) doGroupUpdate(groupID, groupName string, groupVersion int) error {
	return l.db.
		Model(readmodel.ResourceSharingReadModel{}).
		Where("group_key = ? and group_version <= ?", groupID, groupVersion).
		Updates(map[string]interface{}{
			"group_name":    groupName,
			"group_version": groupVersion,
		}).Error
}

func (l *ResourceReadModelHandler) handleUserDiscovered(ctx context.Context, e userdomain.UserDiscovered) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return getOptimisticLocking(l.db, e.SequenceNo).
			Model(&readmodel.ResourceUserNameReadModel{}).
			Create(readmodel.ResourceUserNameReadModel{
				UserKey:  e.AggregateID,
				Username: e.UserInfo.Username,
				Version:  e.SequenceNo,
			}).Error
	})
	g.Go(func() error {
		return l.doUserUpdate(ctx, e.AggregateID, e.UserInfo.Username, e.SequenceNo)
	})
	return g.Wait()
}

func (l *ResourceReadModelHandler) handleUserInfoChanged(ctx context.Context, e userdomain.UserInfoChanged) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return l.db.
			Model(&readmodel.ResourceUserNameReadModel{}).
			Where("user_key = ? and version <= ?", e.AggregateID, e.SequenceNo).
			Updates(map[string]interface{}{
				"username": e.NewUserInfo.Username,
				"version":  e.SequenceNo,
			}).Error
	})
	g.Go(func() error {
		return l.doUserUpdate(ctx, e.AggregateID, e.NewUserInfo.Username, e.SequenceNo)
	})
	return g.Wait()
}

func (l *ResourceReadModelHandler) doUserUpdate(ctx context.Context, userId, username string, version int) error {
	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		return l.db.
			Model(readmodel.DbResourceReadModel{}).
			Where("created_by = ? and created_by_version <= ?", userId, version).
			Updates(map[string]interface{}{
				"created_by_name": username,
			}).Error
	})
	g.Go(func() error {
		return l.db.
			Model(readmodel.DbResourceReadModel{}).
			Where("updated_by = ? and updated_by_version <= ?", userId, version).
			Updates(map[string]interface{}{
				"updated_by_name": username,
			}).Error
	})
	return g.Wait()
}

func (l *ResourceReadModelHandler) handleResourceRegistered(e domain.ResourceRegistered) error {

	resourceKey, err := keys.ParseResourceKey(e.AggregateID)
	if err != nil {
		return err
	}

	var user readmodel.ResourceUserNameReadModel
	qry := l.db.Model(&readmodel.ResourceUserNameReadModel{}).Where("user_key = ?", e.RegisteredBy.String()).Find(&user)
	err = qry.Error
	if err != nil {
		return err
	}

	var username string
	var userVersion = -1
	if qry.RowsAffected > 0 {
		username = user.Username
		userVersion = user.Version
	}

	return getOptimisticLocking(l.db, e.SequenceNo).Create(&readmodel.DbResourceReadModel{
		ResourceReadModelBase: readmodel.ResourceReadModelBase{
			ResourceKey:       resourceKey,
			CreatedBy:         e.RegisteredBy.String(),
			CreatedByName:     username,
			CreatedByVersion:  userVersion,
			CreatedAt:         e.EventTime,
			UpdatedBy:         e.RegisteredBy.String(),
			UpdatedByVersion:  userVersion,
			UpdatedByName:     username,
			UpdatedAt:         e.EventTime,
			GroupSharingCount: 0,
			Version:           e.SequenceNo,
			Owner:             e.RegisteredFor,
		},
		ResourceInfoBase: domain.ResourceInfoBase{
			Name:         e.ResourceInfo.Name,
			Description:  e.ResourceInfo.Description,
			CallType:     e.ResourceInfo.CallType,
			ResourceType: e.ResourceInfo.ResourceType,
		},
	}).Error

}

func (l *ResourceReadModelHandler) handleResourceInfoChanged(e domain.ResourceInfoChanged) error {

	var user readmodel.ResourceUserNameReadModel
	qry := l.db.Model(&readmodel.ResourceUserNameReadModel{}).Where("user_key = ?", e.ChangedBy.String()).Find(&user)
	err := qry.Error
	if err != nil {
		return err
	}

	var username string
	var userVersion = -1
	if qry.RowsAffected > 0 {
		username = user.Username
		userVersion = user.Version
	}

	updates := map[string]interface{}{
		"name":        e.NewResourceInfo.Name,
		"description": e.NewResourceInfo.Description,
		"updated_at":  e.EventTime,
		"updated_by":  e.ChangedBy.String(),
		"version":     e.SequenceNo,
	}
	if username != "" {
		updates["updated_by_name"] = username
	}
	if userVersion != -1 {
		updates["updated_by_version"] = userVersion
	}

	return l.db.
		Model(&readmodel.DbResourceReadModel{}).
		Where("resource_key = ? and version < ?", e.AggregateID, e.SequenceNo).
		Updates(updates).
		Error

}

func (l *ResourceReadModelHandler) handleResourceGroupSharingChanged(ctx context.Context, e domain.ResourceGroupSharingChanged) error {
	resourceKey, err := keys.ParseResourceKey(e.AggregateID)
	if err != nil {
		return err
	}
	g, ctx := errgroup.WithContext(ctx)
	if len(e.RemovedSharings) > 0 {
		g.Go(func() error {
			deleteSql := "resource_key = ? and group_key in ("
			var deleteParams = []interface{}{
				e.AggregateID,
			}
			for i, removedSharing := range e.RemovedSharings {
				deleteSql = deleteSql + "?"
				if i < len(e.RemovedSharings)-1 {
					deleteSql = deleteSql + ","
				}
				deleteParams = append(deleteParams, removedSharing.GroupKey.String())
			}
			deleteSql = deleteSql + ")"
			return l.db.Debug().Where(deleteSql, deleteParams...).Delete(readmodel.ResourceSharingReadModel{}).Error
		})
	}
	g.Go(func() error {
		var toAdd []*readmodel.ResourceSharingReadModel
		for _, addedSharing := range e.AddedSharings {
			var group readmodel.ResourceGroupNameReadModel
			qry := l.db.Model(&readmodel.ResourceGroupNameReadModel{}).Where("group_key = ?", addedSharing.GroupKey.String()).Find(&group)
			if qry.Error != nil {
				return qry.Error
			}
			var groupName string
			var groupVersion = -1
			if qry.RowsAffected > 0 {
				groupName = group.GroupName
				groupVersion = group.Version
			}
			toAdd = append(toAdd, &readmodel.ResourceSharingReadModel{
				ResourceKey:  resourceKey,
				GroupKey:     addedSharing.GroupKey,
				GroupName:    groupName,
				Version:      e.SequenceNo,
				GroupVersion: groupVersion,
			})
		}
		if len(toAdd) > 0 {
			return getOptimisticLocking(l.db, e.SequenceNo).Create(toAdd).Error
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return err
	}
	return l.db.Model(&readmodel.DbResourceReadModel{}).Where("resource_key = ?", resourceKey).Updates(map[string]interface{}{
		"group_sharing_count": len(e.NewResourceSharings),
		"version":             e.SequenceNo,
		"updated_at":          e.EventTime,
	}).Error
}

func (l *ResourceReadModelHandler) handleResourceDeleted(ctx context.Context, e domain.ResourceDeleted) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return l.db.Delete(&readmodel.ResourceSharingReadModel{}, "resource_key = ?", e.AggregateID).Error
	})
	g.Go(func() error {
		return l.db.Delete(&readmodel.DbResourceReadModel{}, "resource_key = ?", e.AggregateID).Error
	})
	return g.Wait()
}

func (l *ResourceReadModelHandler) handleResourceEvaluated(ctx context.Context, e domain.ResourceEvaluated) error {
	resourceKey, err := keys.ParseResourceKey(e.AggregateID)
	if err != nil {
		return err
	}
	g, ctx := errgroup.WithContext(ctx)
	evaluationID := uuid.NewV4().String()
	for _, value := range e.NewEvaluation {
		var rm = &readmodel.ResourceEvaluationReadModel{
			ID:             uuid.NewV4().String(),
			EvaluationID:   evaluationID,
			ResourceKey:    resourceKey,
			EvaluatedBy:    e.EvaluatedBy,
			EvaluatedAt:    e.EventTime,
			DimensionValue: value,
		}
		g.Go(func() error {
			return l.db.Create(rm).Error
		})
	}
	g.Go(func() error {
		return l.db.Model(&readmodel.DbResourceReadModel{}).Where("resource_key = ?", resourceKey).Updates(map[string]interface{}{
			"version":    e.SequenceNo,
			"updated_at": e.EventTime,
		}).Error
	})
	return g.Wait()
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
